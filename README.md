# TASK

- need app which will convert WAV file to MP3 with preinstaller ffmpeg

- module for convert must be reusable module

<br><hr>

## HOW TO RUN
- command line interface: `app -f path`

```bash
./app -f /files/file.wav
```
```bash
./app -f ./files/file.wav
```
```bash
./app -f files/file.wav
```
- ffmpeg can be preinstalled with `brew install ffmeg`


## Requirements:

- Path must consists of 2 elements at least: file path begins with 'files' and file name. File name could be any media file

- Real file must exists

- Real file must be in `output` folder. It means in path provided to app in path part `files` must be replaced by `output`

- in case of any error app must stop and show error message

- ffmpeg output (from any stream, stdout or stderr) must be printed in app output, only convert info

in short:
```
~/pth/app -f files/file.wav <---> ~/pth/app -f ~/pth/output/file.wav
```

output example:
```bash
% ./app/cmd -f /files/alef.wav
2022/10/11 21:49:12 INFO: gonna run FFmpeg
- - -
gonna run FFmpeg
exec: already started
OUT:
Guessed Channel Layout for Input Stream #0.0 : stereo
Input #0, wav, from '/Users/me/dummytask/app/output/alef.wav':
  Metadata:
    encoder         : Lavf59.27.100
  Duration: 00:00:08.51, bitrate: 1411 kb/s
  Stream #0:0: Audio: pcm_s16le ([1][0][0][0] / 0x0001), 44100 Hz, stereo, s16, 1411 kb/s
Stream mapping:
  Stream #0:0 -> #0:0 (pcm_s16le (native) -> mp3 (libmp3lame))
Press [q] to stop, [?] for help
Output #0, mp3, to '/Users/me/dummytask/app/output/alef.mp3':
  Metadata:
    TSSE            : Lavf59.27.100
  Stream #0:0: Audio: mp3, 44100 Hz, stereo, s16p
    Metadata:
      encoder         : Lavc59.37.100 libmp3lame
size=     134kB time=00:00:08.51 bitrate= 128.6kbits/s speed=62.9x    
video:0kB audio:133kB subtitle:0kB other streams:0kB global headers:0kB muxing overhead: 0.185115%


DONE.
```
```bash
% t 3
.
├── app
│  ├── output
│  │  ├── alef.mp3
│  │  └── alef.wav
│  └── cmd
├── cmd
│  └── main.go
├── pkg
│  └── toMP3
│     └── converter.go
├── go.mod
├── Makefile
└── README.md
```