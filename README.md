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
$ vkscraper --login=MyUsername [...]
```

When logging in, **vkscraper** stores the access token in a file called
`./MyUsername.vksession`, which will be reused later the next time `--login` is
given.

Do not delete the session file, logging in is an expensive operation.

### Downloading users and communities

```sh
$ vkscraper --login=MyUsername insidevk
```

By default, all available content is downloaded.

 - Do not download stories: `--no-stories`.
 - Do not download photos: `--no-photos`.
 - Do not download videos: `--no-videos`.
 - Do not download posts (text and attached photos): `--no-posts`.

### Batch file

**vkscraper** can read user profiles and communities from a file. Lines starting
with an `#` or empty lines are considered as comments and ignored. Inline
comments are also ignored.

Given `DataHoarder.txt`:

```sh
# Official VK community.
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
