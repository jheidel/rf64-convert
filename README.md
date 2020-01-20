A utility for converting RIFF WAV files containing IQ data 
[produced by gnuradio](https://wiki.gnuradio.org/index.php/Wav_File_Sink)
into
[RF64](https://en.wikipedia.org/wiki/RF64) WAV files for use with
[SDR-Radio Console](https://www.sdr-radio.com/console).

--


First you'll need a WAV file containing IQ data, for instance from a gnuradio
flow graph like this:

![gnuradio iq source](https://imgur.com/O6oska5.jpg)

In my case, the IQ data comes from a headless Raspberry Pi + AirSpy data logger
running gnuradio.

Then run the conversion:

```
go run main.go --input=[path to gnuradio input WAV] --output=[path to output WAV]
```

In order to get SDR-Radio to display the correct timestamp and center
frequency, I've found the easiest way is to emulate the filenames used by SDR#,
which SDR-Radio will recognize.

Example:

```
--output=SDRSharp_20200119_110726Z_162550000Hz-IQ.wav
```

SDR-Radio Console's own WAV output uses an 'auxi' section within the WAV file
which contains XML-encoded metadata in UTF16 format. `auxi.go` shows the
digging I've done into it so far, but the SDR# path workaround above is
simpler.

73 de KI7QIV!
