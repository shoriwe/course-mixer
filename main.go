package main

import (
	"encoding/json"
	"os"

	"github.com/Woynert/course-mixer/scraper"
	"github.com/Woynert/course-mixer/utils"
)

type Query struct {
	Name string `json:"name"`
	File string `json:"file"`
	Cols uint8  `json:"cols"`
}

type Hour struct {
	Day   string `json:"day"`
	Start int    `json:"start"`
	End   int    `json:"end"`
}

type Course struct {
	Ctg   string `json:"ctg"`
	Level uint8  `json:"level"`
	Title string `json:"title"`
	Nrc   string `json:"nrc"`
	Hours []Hour `json:"hour"`
}

func downloadData() []Course {
	courses, err := scraper.Crawl()
	utils.Fatal(err)
	result := make([]Course, 0, len(courses))
	for _, course := range courses {
		hours := make([]Hour, 0, len(course.Hours))
		for _, hour := range course.Hours {
			h := Hour{
				Day:   hour.Day.String(),
				Start: hour.From.Hour(),
				End:   hour.To.Hour(),
			}
			hours = append(hours, h)
		}
		c := Course{
			Ctg:   course.Faculty,
			Level: course.Level,
			Title: course.Subject,
			Nrc:   course.NRC,
			Hours: hours,
		}
		result = append(result, c)
	}
	return result
}

func main() {
	output := os.Stdout
	if len(os.Args) > 1 {
		var err error
		output, err = os.Create(os.Args[1])
		utils.Fatal(err)
	}
	courses := downloadData()
	encoder := json.NewEncoder(output)
	encoder.SetIndent("", "\t")
	err := encoder.Encode(courses)
	utils.Fatal(err)
}
