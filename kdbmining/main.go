package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/text/unicode/norm"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

type SubjectRef struct {
	CourseID string
	Title string
}

type Subject struct {
	CourseID string
	Title string
	Credit string
	Grade string
	Timetable string
	Books []string
	ClassName []string
	PlanPretopics string
	Keywords []string
	SeeAlsoSubject []*SubjectRef
	Summary string
}

func getCourseID(doc *html.Node) (string, error) {
	id, err := htmlquery.Query(doc, "/html/body/h1[@id='course-title']/span[@id='course']/text()")
	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %s", err)
		return "", err
	}
	return id.Data, err
}

func main() {
	var err error
	doc, _ := htmlquery.Parse(os.Stdin)

	/*
	node, err := htmlquery.Query(doc, "/html/body/h1[@id='course-title']/span[@id='title']/text()")
	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %s", err)
		return
	}
	fmt.Println(node.Data)
	 */
	subject := new(Subject)

	subject.CourseID, err = getCourseID(doc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %s", err)
		return
	}
	subject.CourseID = strings.TrimSpace(subject.CourseID)

	title, err := htmlquery.Query(doc, "/html/body/h1[@id='course-title']/span[@id='title']/text()")
	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %s", err)
		return
	}
	subject.Title = strings.TrimSpace(title.Data)

	credit, err := htmlquery.Query(doc, "/html/body/div[@id='credit-grade-assignments']/p/span[@id='credit']/text()")
	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %s", err)
		return
	}
	subject.Credit = strings.TrimSpace(
		strings.TrimRight(credit.Data, ","),
	)

	grade, err := htmlquery.Query(doc, "/html/body/div[@id='credit-grade-assignments']/p/span[@id='grade']/text()")
	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %s", err)
		return
	}
	subject.Grade = strings.TrimSpace(
		strings.TrimRight(grade.Data, ","),
	)

	timetable, err := htmlquery.Query(doc, "/html/body/div[@id='credit-grade-assignments']/p/span[@id='timetable']/text()")
	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %s", err)
		return
	}
	subject.Timetable = strings.TrimSpace(
		timetable.Data,
	)

	books, err := htmlquery.QueryAll(doc, "//h2[text()='教材・参考文献・配付資料等']/following-sibling::div[@id='topics']//div/table//th/text()")
	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %s", err)
		return
	}
	for _, book := range books {
		subject.Books = append(
			subject.Books,
			strings.TrimSpace(
				book.Data,
			),
		)
	}

	plan, err := htmlquery.Query(doc, "//h2[text()='授業計画']/following-sibling::div[@id='topics']")
	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %s", err)
		return
	}

	classes, err := htmlquery.QueryAll(plan, "//div/table//td[not(contains(@nowrap, 'nowrap'))]/text()")
	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %s", err)
		return
	}
	for _, class := range classes {
		subject.ClassName = append(
			subject.ClassName,
			strings.TrimSpace(
				class.Data,
			),
		)
	}

	plan_topics, err := htmlquery.QueryAll(plan, "//p[@id='pretopics']/text()")
	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %s", err)
		return
	}
	for _, topic := range plan_topics {
		subject.PlanPretopics = subject.PlanPretopics + strings.TrimSpace(
			topic.Data,
		)
	}

	related_subjects_elem, err := htmlquery.QueryAll(doc, "/html/body/div/h2[text()='他の授業科目との関連']/parent::node()/following-sibling::div[@id='topic-assignments']//table/tbody/tr")
	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %s", err)
		return
	}
	for _, subject_elem := range related_subjects_elem {
		course_id, _ := htmlquery.Query(subject_elem, "//th/text()")
		course_name, _ := htmlquery.Query(subject_elem, "//td/text()")

		ref := new(SubjectRef)
		ref.CourseID = strings.TrimSpace(course_name.Data)
		ref.Title = strings.TrimSpace(course_id.Data)
		subject.SeeAlsoSubject = append(subject.SeeAlsoSubject, ref)
	}

	keyword, err := htmlquery.Query(doc, "//h2[text()='キーワード']//following-sibling::p[@id='style']/text()")
	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %s", err)
		return
	}
	if keyword != nil {
		normalized_keyword := norm.NFKC.String(keyword.Data)
		for _, value := range strings.Split(normalized_keyword, ",") {
			keyword_str := strings.TrimSpace(
				value,
			)
			subject.Keywords = append(subject.Keywords, keyword_str)
		}
	}

	summary, err := htmlquery.Query(doc, "//p[@id='summary-contents']/text()")
	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %s", err)
		return
	}
	subject.Summary = summary.Data
	fmt.Println(subject)
}
