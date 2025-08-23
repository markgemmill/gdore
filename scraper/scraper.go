package scraper

import (
	"fmt"
	"gdore/environ"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func OpenMSEdgeBrowser(env environ.Environ, logger *Logger) string {
	logger.Log("launching MS Edge browser...")
	return launcher.New().Bin(env.BrowserApp).Headless(false).MustLaunch()
}

func OpenDefaultBrowser(env environ.Environ, logger *Logger) string {
	logger.Log("launching default Chromium browser...")
	return launcher.New().Headless(false).MustLaunch()
}

func OpenPortal(env environ.Environ, logger *Logger) *rod.Page {
	u := ""
	if env.BrowserApp == "" {
		u = OpenDefaultBrowser(env, logger)
	} else {
		u = OpenMSEdgeBrowser(env, logger)
	}
	logger.Log("opening sobeys portal...")
	page := rod.New().ControlURL(u).MustConnect().MustPage(SOBEYS_PORTAL)
	return page
}

func Login(page *rod.Page, userId string, password string, logger *Logger) {
	// On the first page is the login
	// A successful login navigates to the management page
	logger.Log("logging in...")
	messanger.SendMsg("logging in...")
	page.MustWaitStable()
	page.MustElement(LOGIN_USER_ID).MustInput(userId)
	page.MustElement(LOGIN_USER_PW_ID).MustInput(password)
	page.MustElement(LOGIN_LANGUAGE_ID).MustSelect("English")
	page.MustElement(LOGIN_USER_SUBMIT_ID).MustClick()
}

func LoginSuccessful(page *rod.Page, logger *Logger) bool {
	page.MustWaitStable()

	LOGIN_MSG_ID := "#MainContent_lblMessage"

	failedLogin, err := page.Timeout(2 * time.Second).Element(LOGIN_MSG_ID)

	logger.Logf("login err: %s", err)
	logger.Logf("login err msg element: %s", failedLogin)

	if failedLogin != nil {
		// login failed
		failedLogin.CancelTimeout()

		errMsg, err := failedLogin.Text()
		if err != nil {
			logger.Logf("failed to collect error message text: %s", err)
		}

		logger.Log(errMsg)

		messanger.SendError("login failed", fmt.Errorf("login failed"))
		return false
	}

	// For success we should fail to find the element...
	messanger.SendMsg("login successful")
	return true
}

func GoToManagement(page *rod.Page) {
	// We are on the management page, and click on
	// the Cost and Deals Management link.
	// This will take us to the deal search page
	page.MustWaitStable()
	page.MustElementR("a", "Cost and Deals Management").MustClick()
}

func DealSearch(page *rod.Page, region string, dealNumber string, logger *Logger) {
	// We are on the deal search window.
	// Here we enter the individual deal number
	// which will take us to the search resuls page.
	messanger.SendMsg("looking up %s", dealNumber)
	logger.Logf("look up %s", dealNumber)
	logger.Logf("input %s", dealNumber)
	page.MustElement(DEAL_NO_ID).MustInput(dealNumber)
	logger.Logf("input region %s", REGION_NO)
	page.MustElement(REGION_NO_ID).MustSelect(REGION_NO)

	// set to vendor (not broker)
	// not all searches have a vendor/broker option
	logger.Log("set to vendor")
	el, err := page.Timeout(2 * time.Second).Element(VENDOR_RADIO_ID)
	if err == nil {
		el.MustClick()
	}

	logger.Log("click search")
	page.MustElement(SEARCH_ID).MustClick()
}

func SearchResults(page *rod.Page, dealNumber string, logger *Logger) (*rod.Element, error) {
	logger.Logf("waiting for search page on %s", dealNumber)
	page.MustWaitStable()
	el, err := page.Timeout(3*time.Second).ElementR("a", dealNumber)
	if err != nil {
		logger.Log("search results did not load")
		return nil, err
	}

	logger.Log("found deal link on search page")
	return el, nil
}

func GoToLookup(page *rod.Page, logger *Logger) {
	page.MustElement(LOOKUP_BTN_ID).MustClick()
	page.MustWaitStable()
	logger.Log("back to deal lookup")
}

var docIdRx regexp.Regexp = *regexp.MustCompile(`^\d+$`)

func ParseDocumentIdInput(documents string) []string {
	docs := []string{}
	for _, line := range slices.Sorted(strings.Lines(documents)) {
		doc := strings.Trim(line, "\r\n ")
		if doc == "" {
			continue
		}
		if docIdRx.MatchString(doc) {
			docs = append(docs, doc)
		}
	}

	return docs
}

func RunScraper(userId, password, region, documentIds string, env environ.Environ) {
	logger, err := InitLogger(env.LogDir)
	if err != nil {
		messanger.SendError("error creating logger", err)
		return
	}
	defer logger.Close()

	page := OpenPortal(env, logger)

	Login(page, userId, password, logger)
	if !LoginSuccessful(page, logger) {
		return
	}

	GoToManagement(page)
	Pause(1)

	docs := ParseDocumentIdInput(documentIds)

	docCount := len(docs)
	downloadedDocCount := 0
	erroredDocCount := 0

	if docCount == 0 {
		logger.Log("no valid documents")
		messanger.SendDone("no valid docments", "")
		return
	}

	output := NewExcelDocument(env.OutputDir)

	for index, docId := range docs {
		logger.Logf("%d: %s", index, docId)

		DealSearch(page, region, docId, logger)
		// Pause(2)

		link, err := SearchResults(page, docId, logger)
		if err != nil {
			logger.Logf("document %s not found", docId)
			messanger.SendMsg("document %s not found", docId)
			erroredDocCount += 1
			Pause(2)

			ClearErrorDialog(page)
			fmt.Println("dialog cleared")

			// GoToLookup(page, logger)

			continue
		}

		logger.Logf("clicking deal %s link", docId)
		link.MustClick()
		page.MustWaitStable()

		logger.Logf("start parsing deal %s data", docId)
		// extract document data and print
		document := NewDocument(docId)
		ParseDocumentInfo(page, document, logger)
		ParseDocumentDates(page, document)
		ParseLogistics(page, document)
		ParseArticleDetails(page, document)

		PrintDocument(document)
		output.AddDocument(document)
		downloadedDocCount += 1
		messanger.SendMsg("%d of %d documents downloaded", downloadedDocCount, docCount)
		Pause(1)

		GoToLookup(page, logger)
	}

	completionMsg := fmt.Sprintf("Completed. %d of %d documents downloaded.", downloadedDocCount, docCount)
	logger.Log(completionMsg)

	if downloadedDocCount > 0 {
		logger.Log(output.path.String())
		err := output.Save()
		if err != nil {
			completionMsg = fmt.Sprint("%s Failed to save results!", completionMsg)
			logger.Logf("error: %s", err)
		}
	}

	messanger.SendDone(completionMsg, output.path.String())
	page.Close()
}
