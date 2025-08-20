package scraper

var (
	// input variables
	SOBEYS_PORTAL = "https://partner.sobeys.com/Login"
	REGION_NO     = "Quebec"
	// loging page selectors
	LOGIN_USER_ID        = "#MainContent_txtUsername"
	LOGIN_USER_PW_ID     = "#MainContent_txtPassword"
	LOGIN_LANGUAGE_ID    = "#MainContent_selLanguage"
	LOGIN_USER_SUBMIT_ID = "#MainContent_btnSubmit"
	// deal search
	DEAL_NO_ID      = "#dlu_deal_num"
	REGION_NO_ID    = "#dlu_region2"
	VENDOR_RADIO_ID = `input[type="radio"][value="v"]`
	SEARCH_ID       = "#dlu_btn_search_doc_num"
	LOOKUP_BTN_ID   = "#dve_btn_back_to_lookup"

	// deal page selectors
	DEAL_SELECTOR            = `a[text="%s"]`
	ARTICLES_TABLE           = "#dve_articles_table"
	DOC_TYPE_INPUT_ID        = "#dve_doc_type_descr"
	DOC_SOURCE_TYPE_INPUT_ID = "#dve_source_type_desc"
	DOC_STATUS_INPUT_ID      = "#dve_doc_status_desc"
	PROMO_EFF_FROM_ID        = "#dve_store_promo_eff_from_mirror"
	PROMO_EFF_TO_ID          = "#dve_store_promo_eff_to_mirror"

	VENDOR_INPUT_ID = "#dve_vendor"
	REGION_INPUT_ID = "#dve_region_desc"
)
