package scraper

import (
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

// • Document No -- DONE
// • Document Type -- DONE
// • Warehouse/DSD -- DONE
// • Promo Effective Dates -- DONE
// • Articles:
//   - GTIN - DONE
//   - Number - DONE
//   - Description - DONE
//   - Article Count/ReCount
//
// • Banner Descriptions
// • Loyalty Buy / Get fields
func ParseDocumentInfo(page *rod.Page, doc *Document, logger *Logger) {
	inputEl := page.MustElement(DOC_TYPE_INPUT_ID)
	doc.DocType = ExtractInputValue(inputEl)
	logger.Logf("%s = '%s'", DOC_TYPE_INPUT_ID, doc.DocType)

	inputEl = page.MustElement(DOC_SOURCE_TYPE_INPUT_ID)
	doc.SourceType = ExtractInputValue(inputEl)
	logger.Logf("%s = '%s'", DOC_SOURCE_TYPE_INPUT_ID, doc.SourceType)

	inputEl = page.MustElement(DOC_STATUS_INPUT_ID)
	doc.Status = ExtractInputValue(inputEl)
	logger.Logf("%s = '%s'", DOC_STATUS_INPUT_ID, doc.Status)

	inputEl = page.MustElement(VENDOR_INPUT_ID)
	doc.Vendor = ExtractInputValue(inputEl)
	logger.Logf("%s = '%s'", VENDOR_INPUT_ID, doc.Vendor)

	inputEl = page.MustElement(REGION_INPUT_ID)
	doc.Region = ExtractInputValue(inputEl)
	logger.Logf("%s = '%s'", REGION_INPUT_ID, doc.Vendor)
}

func ParseDocumentDates(page *rod.Page, doc *Document) {
	// navigate to the articles tab
	page.MustElementR("a", "Deals").MustClick()
	page.MustWaitStable()

	fmt.Println("Extracting Effective From Date")
	inputEl := page.MustElement(PROMO_EFF_FROM_ID)
	doc.EffectiveFrom = ExtractInputValue(inputEl)

	fmt.Println("Extracting Effective From Date")
	inputEl = page.MustElement(PROMO_EFF_TO_ID)
	doc.EffectiveTo = ExtractInputValue(inputEl)
}

type DealType struct {
	ColumnIndex int
	TypeName    string
}

func ParseDealType(header *rod.Element) DealType {
	row := header.MustElement("tr")
	cells := row.MustElements("th")
	for index, cell := range cells {
		if index < 14 || index > 19 {
			continue
		}
		visible, _ := cell.Visible()
		if visible {
			switch index {
			case 14:
				return DealType{ColumnIndex: index, TypeName: "Off Invoice"}
			case 15:
				return DealType{ColumnIndex: index, TypeName: "By Cheque"}
			case 16:
				return DealType{ColumnIndex: index, TypeName: "Count / Recount"}
			case 17:
				return DealType{ColumnIndex: index, TypeName: "Scan Funding"}
			case 18:
				return DealType{ColumnIndex: index, TypeName: "New Cost"}
			case 19:
				return DealType{ColumnIndex: index, TypeName: "Quantity"}
			}
		}
	}
	return DealType{ColumnIndex: 14, TypeName: "No Deal Type Found"}
}

func ParseArticleDetails(page *rod.Page, doc *Document) {
	// navigate to the articles tab
	page.MustElementR("a", "Articles").MustClick()
	table := page.MustElement(ARTICLES_TABLE)

	header := table.MustElement("thead")
	dealType := ParseDealType(header)

	body := table.MustElement("tbody")

	rows := body.MustElements("tr")
	for _, row := range rows {
		article := Article{}
		// <td class="article-family-field">C307932E1</td>
		article.Family = GetElementText(row, ".article-family-field")
		// <td class="article-desc-field">GoBio Vegetable Cubes 15X66G</td>
		article.Description = GetElementText(row, ".article-desc-field")
		// <td class="article-pack-field">15</td>
		article.Pack = GetElementText(row, ".article-pack-field")
		// <td class="article-size-field">66.000 G</td>
		article.Size = GetElementText(row, ".article-size-field")

		cells := row.MustElements("td")
		if len(cells) > 2 {
			// <input type="text" class="article-gtin-field" maxlength="18" readonly="">
			article.Gtin = GetInputText(cells[1], "input.article-gtin-field")

			// <input type="text" class="article-num-field" maxlength="18" readonly="">
			article.Number = GetInputText(cells[2], "input.article-num-field")

			article.Deal = GetInputText(cells[dealType.ColumnIndex], "input.article-deal-field")
			article.DealType = dealType.TypeName
		}

		doc.AddArticle(article)
	}
}

func ParseLogistics(page *rod.Page, doc *Document) {
	// navigate to the articles tab
	page.MustElementR("a", "Logistics").MustClick()
	table := page.MustElement("#dve_banners_table > tbody")

	rows := table.MustElements("tr")
	if len(rows) == 0 {
		fmt.Println("No table rows found for #dve_banners_table")
		return
	}
	for _, row := range rows {
		banner := Banner{}
		cell := row.MustElement("td.banner-checkbox-column")
		banner.Active = GetInputBool(cell, "input[type='checkbox']")

		cell = row.MustElement("td.banner-id-column")
		banner.Id = GetText(cell)

		cell = row.MustElement("td.banner-desc-column")
		banner.Name = GetText(cell)
		fmt.Printf("Banner: %v\n", banner)
		doc.AddBanner(banner)
	}
}

// Check for a lookup error message and clear it if found
func ClearErrorDialog(page *rod.Page) {
	dialogs, err := page.Timeout(3 * time.Second).Elements(".ui-dialog")
	if err != nil {
		return
	}
	var errDialog *rod.Element
	for _, dialog := range dialogs {
		// aria-labelledby="ui-dialog-title-error_dialog"
		aria, err := dialog.Attribute("aria-labelledby")
		if err != nil {
			continue
		}
		if *aria == "ui-dialog-title-error_dialog" {
			errDialog = dialog
		}
	}

	if errDialog == nil {
		return
	}

	visible, err := errDialog.Visible()
	if err != nil {
		return
	}

	if !visible {
		return
	}
	fmt.Println("found error dialog and clearing...")
	btn, err := errDialog.Timeout(5 * time.Second).Element("button")
	if err != nil {
		fmt.Println("Okay button not found")
		return
	}
	fmt.Println(btn)
	btn.Timeout(2*time.Second).Click(proto.InputMouseButtonLeft, 1)
	fmt.Println(btn)
}

// <input type="checkbox" class="article-checkbox-field" readonly="" disabled="">
// <input type="text" class="article-gtin-field" maxlength="18" readonly="">
// <input type="text" class="article-num-field" maxlength="18" readonly="">
// <td class="article-family-field">C307932E1</td>
// <td class="article-desc-field">GoBio Vegetable Cubes 15X66G</td>
// <input type="checkbox" class="exclart-checkbox-field" readonly="" disabled="" value="">
// <td class="article-pack-field">15</td>
// <td class="article-size-field">66.000 G</td>
// <td class="article-uom-field"></td>
// <td class="article-uom-field">CS</td>
// <td class="article-uom-field">CS</td>
// <td class="article-pack-field">0</td>
// <td class="article-uom-field">EA</td>
// <td class="article-uom-field">40.37</td>
// <td>
// <input type="text" class="article-deal-field" readonly="">
// </td>
// <td style="display: none;">
// <input type="text" class="article-deal-field" maxlength="12" readonly="">
// </td><td style="display: none;">
// <input type="text" class="article-deal-field" maxlength="12" readonly="">
// <input type="text" class="article-deal-field" maxlength="12" readonly="">
// <input type="text" class="article-deal-field" maxlength="12" readonly="">
// <input type="text" class="article-deal-field" maxlength="14" readonly="">
// <td class="article-sugg-retail-col" style="display: none;">
// <input type="text" class="article-sugg-retail-field" maxlength="14" readonly="" style="display: inline;">
// <input type="text" class="article-sugg-retail-field" maxlength="12" readonly="" style="display: inline;">
// <select disabled="" style="display: inline;">
//   <option></option>
//   <option value="CS">CS</option>
//   <option value="CAR">CAR</option>
//   <option value="EA">EA</option>
//   <option value="KG">KG</option>
// </select>
// <input type="text" class="article-comments-field" maxlength="40" readonly="">
// <td class="articleBannerList" style="border-bottom: 1px solid rgb(153, 153, 153);">
// <span style="display:inline-block;width:30px;"> </span>
// <span style="border-right:#999 1px solid;"></span>
// <span style="display:inline-block;width:30px;">81</span>
// <span style="border-right:#999 1px solid;"></span><span style="display:inline-block;width:30px;"> </span><span style="border-right:#999 1px solid;"></span><span style="display:inline-block;width:30px;"> </span><span style="border-right:#999 1px solid;"></span><span style="display:inline-block;width:30px;"> </span><span style="border-right:#999 1px solid;"></span><span style="display:inline-block;width:30px;"> </span><span style="border-right:#999 1px solid;"></span><span style="display:inline-block;width:30px;"> </span><span style="border-right:#999 1px solid;"></span><span style="display:inline-block;width:30px;"> </span><span style="border-right:#999 1px solid;"></span><span style="display:inline-block;width:30px;"> </span><span style="border-right:#999 1px solid;"></span><span style="display:inline-block;width:30px;"> </span><span style="border-right:#999 1px solid;"></span><span style="display:inline-block;width:30px;"> </span><span style="border-right:#999 1px solid;"></span><span style="display:inline-block;width:30px;"> </span><span style="border-right:#999 1px solid;"></span><span style="display:inline-block;width:30px;"> </span><span style="border-right:#999 1px solid;"></span><span style="display:inline-block;width:30px;"> </span><span style="border-right:#999 1px solid;"></span><span style="display:inline-block;width:30px;"> </span><span style="border-right:#999 1px solid;"></span><span style="display:inline-block;width:30px;"> </span><span style="border-right:#999 1px solid;"></span><span style="display:inline-block;width:30px;"> </span><span style="border-right:#999 1px solid;"></span><span style="display:inline-block;width:30px;"> </span><span style="border-right:#999 1px solid;"></span>
// <span style="display:inline-block;width:30px;"> </span></td>
// <td class="article-price-family-field" style="display: none;">R307932E1</td></tr>
