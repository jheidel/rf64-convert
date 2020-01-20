package convert

import (
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

var (
	// SDRConsole WAV files contain an 'auxi' section with the following metadata.
	TestAuxi = "<?xml version=\"1.0\"?>\r\n<SDR-XML-Root xml:lang=\"EN\" Description=\"Saved recording data\" Created=\"20-Jan-2020 01:50\"><Definition CurrentTimeUTC=\"20-01-2020 01:50:24\" Filename=\"19-Jan-2020 175024.012 162.550MHz.wav\" FirstFile=\"19-Jan-2020 175024.012 162.550MHz.wav\" Folder=\"F:\\vbox2\" InternalTag=\"5E25-0760-000C\" PreviousFile=\"\" RadioModel=\"RTL Dongle - R820T\" RadioSerial=\"\" SoftwareName=\"SDR Console\" SoftwareVersion=\"Version 3.0.18 build 1740\" UTC=\"20-01-2020 01:50:24\" XMLLevel=\"XMLLevel003\" CreatedBy=\"Jeff on JH-FLANDRE\" TimeZoneStatus=\"1\" TimeZoneInfo=\"4AEAAFAAYQBjAGkAZgBpAGMAIABTAHQAYQBuAGQAYQByAGQAIABUAGkAbQBlAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAsAAAABAAIAAAAAAAAAAAAAAFAAYQBjAGkAZgBpAGMAIABEAGEAeQBsAGkAZwBoAHQAIABUAGkAbQBlAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAMAAAACAAIAAAAAAAAAxP///w==\" DualMode=\"0\" Sequence=\"0\" ADFrequency=\"0\" BitsPerSample=\"16\" BytesPerSecond=\"10000000\" RadioCenterFreq=\"162550000\" SampleRate=\"2500000\" UTCSeconds=\"1579485024\"></Definition></SDR-XML-Root>\r\n\x00"
)

func DecodeUTF16(b []byte) (string, error) {
	r, _, err := transform.Bytes(unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder(), b)
	return string(r), err
}

func EncodeUTF16(s string) ([]byte, error) {
	r, _, err := transform.Bytes(unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder(), []byte(s))
	return r, err
}
