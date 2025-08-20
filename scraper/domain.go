package scraper

import (
	"fmt"
	"strings"
)

type Document struct {
	Number        string
	DocType       string
	SourceType    string
	Status        string
	Vendor        string
	Region        string
	EffectiveFrom string
	EffectiveTo   string
	Articles      []Article
	Banners       []Banner
}

func NewDocument(documentId string) *Document {
	return &Document{
		Number:   documentId,
		Articles: []Article{},
		Banners:  []Banner{},
	}
}

func PrintDocument(doc *Document) {
	fmt.Printf("Doc:  %s\n", doc.Number)
	fmt.Printf("Type: %s\n", doc.DocType)
	fmt.Printf("Src:  %s\n", doc.SourceType)
	fmt.Printf("From: %s\n", doc.EffectiveFrom)
	fmt.Printf("To:   %s\n", doc.EffectiveTo)
	fmt.Println("Articles:")
	for index, art := range doc.Articles {
		fmt.Printf(" %d: %s %s\n", index, art.Number, art.Description)
	}
}

func (d *Document) AddArticle(article Article) {
	d.Articles = append(d.Articles, article)
}

func (d *Document) AddBanner(banner Banner) {
	d.Banners = append(d.Banners, banner)
}

func (d *Document) Headers() []string {
	return []string{
		"Document Number",
		"Document Type",
		"Source Type",
		"Status",
		"Vendor",
		"Region",
		"Effective From",
		"Effective To",
	}
}

func (d *Document) Values() []string {
	return []string{
		d.Number,
		d.DocType,
		d.SourceType,
		d.Status,
		d.Vendor,
		d.Region,
		d.EffectiveFrom,
		d.EffectiveTo,
	}
}

func (d *Document) String() string {
	return strings.Join(d.Values(), ",")
}

func (d *Document) Csv() string {
	rows := []string{}
	for _, article := range d.Articles {
		rows = append(rows, d.String()+","+article.String())
	}
	return strings.Join(rows, "\n")
}

func (d *Document) Table() [][]string {
	table := [][]string{}
	for _, article := range d.Articles {
		for _, banner := range d.Banners {
			if !banner.Active {
				continue
			}
			row := []string{}
			row = append(row, d.Values()...)
			row = append(row, banner.Values()...)
			row = append(row, article.Values()...)
			table = append(table, row)
		}
	}
	return table
}

type Article struct {
	Gtin        string
	Number      string
	Family      string
	Description string
	Pack        string
	Size        string
	Deal        string
	DealType    string
}

func (a Article) Values() []string {
	return []string{
		a.Gtin,
		a.Number,
		a.Family,
		a.Description,
		a.Pack,
		a.Size,
		a.Deal,
		a.DealType,
	}
}

func (a Article) Headers() []string {
	return []string{
		"GTIN",
		"Number",
		"Family",
		"Description",
		"Pack",
		"Size",
		"Deal",
		"Deal Type",
	}
}

func (a Article) String() string {
	return strings.Join(a.Values(), ",")
}

type Banner struct {
	Active bool
	Id     string
	Name   string
}

func (b Banner) Headers() []string {
	return []string{
		"Banner Id",
		"Banner",
	}
}

func (b Banner) Values() []string {
	return []string{
		b.Id,
		b.Name,
	}
}
