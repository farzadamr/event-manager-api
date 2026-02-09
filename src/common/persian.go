package common

import (
	"fmt"
	"log"
	"regexp"
	"time"

	persian "github.com/yaa110/go-persian-calendar"
)

const iranianMobileNumberPattern string = `^09(1[0-9]|2[0-2]|3[0-9]|9[0-9])[0-9]{7}$`

func IranianMobileNumberValidate(mobileNumber string) bool {
	res, err := regexp.MatchString(iranianMobileNumberPattern, mobileNumber)
	if err != nil {
		log.Print(err.Error())
	}
	return res
}

func ToShamsiString(t time.Time) string {
	p := persian.New(t)
	return fmt.Sprintf(
		"%04d/%02d/%02d",
		p.Year(), p.Month(), p.Day())
}
