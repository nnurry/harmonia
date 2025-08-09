package utils

import (
	"bytes"
	"fmt"
	"time"
)

func GenerateUniqueTimestamp() string {
	stringBuffer := bytes.NewBufferString("")

	now := time.Now()
	fmt.Fprintf(stringBuffer, "/%v", now.Format("20060102150405"))
	fmt.Fprintf(stringBuffer, "%03d", now.Nanosecond()/1_000_000)

	return stringBuffer.String()
}
