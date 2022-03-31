# vkscraper

Download posts, photos, and videos along with their captions and other metadata
from ВКонтакте.

[![CI](https://github.com/kandayo/vkscraper/actions/workflows/ci.yml/badge.svg)](https://github.com/kandayo/vkscraper/actions/workflows/ci.yml)

## Features

 - Download posts, photos and videos from users profiles and communities.
 - Built with archivists in mind: incremental sync.
 - JSON metadata is stored alongside the media.
 - Batch file, cronjob friendly.

## Usage

### Authentication

```
$ vkscraper --login=MyUsername [--password=MyPassword] [...]
```

When logging in, **vkscraper** stores the access token in a file called
`./MyUsername.vksession`, which will be reused later the next time `--login` is
given. Do not delete the session file, logging in is an expensive operation.

### Downloading users and communities

```sh
$ vkscraper --login=MyUsername insidevk
```

By default, all available content will be downloaded.

#### Archive options

 - `--no-photos`
 - `--no-posts`
 - `--no-stories`
 - `--no-tagged-photos`
 - `--no-videos`

#### Directory structure

```
.
├── photos
│   ├── 2022-03-02T10:20:28Z (34894549853).jpg
│   └── 2022-03-02T10:20:28Z (34894549853).json
├── posts
│   ├── 2022-03-02T10:20:15Z (14312312311).json
│   └── 2022-03-02T10:20:15Z (14312312311).txt
├── stories
│   ├── 2022-01-01T10:20:15Z (54894549852).jpg
│   ├── 2022-01-01T10:20:15Z (54894549852).json
│   ├── 2022-01-01T10:20:28Z (54894549853).mp4
│   └── 2022-01-01T10:20:28Z (54894549853).json
├── tagged_photos
│   ├── 2022-03-02T10:20:28Z (84893123153).jpg
│   └── 2022-03-02T10:20:28Z (84893123153).json
├── videos
│   ├── 2022-01-01T10:20:28Z (64812111853).mp4
│   └── 2022-01-01T10:20:28Z (64812111853).json
├── meta.json
└── id
```

### Batch file

**vkscraper** can read user profiles and communities from a file. Lines starting
with an `#` or empty lines are considered as comments and ignored. Inline
comments are also ignored.

Given `DataHoarder.txt`:

```sh
# Official VK community
insidevk

klavdiacoca # Inline comment 1; Клава Кока; profile
klavacoca   # Inline comment 2; Клава Кока; community
```

Download stories, photos, videos and posts:

```sh
$ vkscraper --login=MyUsername --batch-file=DataHoarder.txt
```

### Fast update

For each target, stop when encountering the first already-downloaded resource.
This option is recommended when you use **vkscraper** to update your personal
archive.

This option was taken from Instaloader.

```sh
$ vkscraper --login=MyUsername --fast-update insidevk
```

## API

**vkscraper** is not meant to be a comprehensive API client for VK. The
functions and structures used internally may be imported as a library.

```go
import "github.com/kandayo/vkscraper/pkg/vk"

vk := vk.NewClient()

// Login with a username and password.
vk.Login("username", "password")
// Or set an access token.
vk.SetAccessToken("token")

// Find the user or community id.
user, err := vk.Utils.ResolveScreenName("klavacoca")

// Retrieve the user stories feed.
stories, err := vk.Stories.Get(user.ID)

perPage := 100
initialOffset := 0

// Retrieve the user photos.
stories, err := vk.Photos.GetAll(user.ID, perPage, initialOffset)

// Retrieve the user videos.
stories, err := vk.Videos.Get(user.ID, perPage, initialOffset)
```
