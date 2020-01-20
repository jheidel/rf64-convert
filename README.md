# RIFF to RF64 WAV Converter

[![Build Status](https://travis-ci.org/jheidel/rf64-convert.svg?branch=master)](https://travis-ci.org/jheidel/rf64-convert)

A utility for converting RIFF WAV files containing IQ data 
[produced by gnuradio](https://wiki.gnuradio.org/index.php/Wav_File_Sink)
into
[RF64](https://en.wikipedia.org/wiki/RF64) WAV files
([spec](https://tech.ebu.ch/docs/tech/tech3306v1_0.pdf)) for use with
[SDR-Radio Console](https://www.sdr-radio.com/console).

--

The standard RIFF WAV format has a file size limitation of 2 GiB because the
size of the `data` chunk is stored within the file as a 32-bit integer. SDR IQ
recordings can easily exceed this restriction due to high sample rates.

The gnuradio WAV Sink will generate RIFF WAV files with a standard header,
followed by an unrestricted data payload. However, the size information in the
final file may not be correct due to integer overflow, and some tools will fail
to load the file due to the incorrect length information. From my experiments,
SDR-Radio Console will end up showing an arbitrary length for these recordings
and will not play the entire file.

Fortunately, SDR-Radio Console supports the RF64 WAV format which handles WAV
files up to 16 EiB. This tool implements conversion to the RF64 WAV file format
by rewriting the WAV file header with a RF64 header and corrected size
information.

--

First you'll need a WAV file containing IQ data, for instance from a gnuradio
flow graph like this:

![gnuradio iq source](https://imgur.com/O6oska5.jpg)

In my case, the IQ data comes from a headless Raspberry Pi 3 + AirSpy data
logger running gnuradio.

Then run the conversion:

```bash
# Install golang
sudo apt install golang

# Fetch repo
git clone https://github.com/jheidel/rf64-convert.git ~/go/src/rf64-convert
cd ~/go/src/rf64-convert

# Build
go build

# Convert
./rf64-convert --input=[path to gnuradio input WAV] --output=[path to output WAV]
```

The input is assumed to be a RIFF WAV file which contains a standard header,
followed by a data chunk which continues to the end of the file. The converter
will rewrite the file with a RF64 header containing the correct file-size
values.

The converter supports using a unix
[pipe](http://man7.org/linux/man-pages/man2/pipe.2.html) as an input for use as
part of a streaming conversion. The output must support seeking so the ds64
header can be updated at the end of the conversion.

In order to get SDR-Radio to display the correct timestamp and center
frequency, I've found the easiest way is to replicate the filenames used by
SDR#, which SDR-Radio will recognize.

Example:

```
--output=SDRSharp_20200119_110726Z_162550000Hz-IQ.wav
```

![sdr radio recordings pane](https://imgur.com/FPELcVH.jpg)

SDR-Radio Console's own WAV output uses a special `auxi` chunk within the WAV
file which contains XML metadata as a UTF16-encoded string. `auxi.go` shows the
digging I've done into it so far and some sample data, but the SDR# path
workaround above is simpler.

73 de KI7QIV!
