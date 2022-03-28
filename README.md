# vkscraper

> **vkscraper** is in its very early alpha stage. The public interface is
> subject to changes and will likely be entirely refactored.

Download posts, photos, and videos along with their captions and other metadata
from ВКонтакте.

[![CI](https://github.com/kandayo/vkscraper/actions/workflows/ci.yml/badge.svg)](https://github.com/kandayo/vkscraper/actions/workflows/ci.yml)

## Features

 - Download posts, photos and videos from users profiles and communities.
 - Built with archivists in mind: incremental sync.
 - JSON metadata is stored alongside the media.
 - Batch file, cronjob friendly.

## Roadmap

 - Better error handling; retry on connection error.
 - Extract comments from posts and media.
 - Schedule integration tests on Github Actions.

## Usage

### Getting an Access Token

All following examples will assume that you pass an `--access-token` or set
`VKSCRAPER_ACCESS_TOKEN` environment variable.

See: https://vk.com/dev/access_token.

### Downloading users and communities

```sh
$ vkscraper --stories --no-photos insidevk
```

By default, only photos are downloaded.

 - Download stories: `--stories`.
 - Do not download photos: `--no-photos`.
 - Download videos: `--videos`.
 - Download posts (text and attached photos): `--posts`.

### Batch file

vkscraper can read user profiles and communities from a file. Lines starting
with an `#` or empty lines are considered as comments and ignored. Inline
comments are also ignored.

Given `batch.txt`:

```sh
# ВКонтакте official account
insidevk

99999 # inline comment 1
88888 # inline comment 2
```

Download posts, photos, videos and stories:

```sh
$ vkscraper --batch-file batch.txt
```

### Fast update

For each target, stop when encountering the first already-downloaded resource.
This option is recommended when you use vkscraper to update your personal
archive.

This option was taken from Instaloader.

```sh
$ vkscraper --stories --videos --posts --fast-update insidevk

Retrieving stories from profile "insidevk".
[1/1] insidevk/stories/111111111.jpg exists
Retrieving photos from profile "insidevk".
[1/745] insidevk/photos/222222222.jpg exists
Retrieving videos from profile "insidevk".
[1/1] insidevk/videos/333333333.mp4 exists
Retrieving posts from profile "insidevk".
[1/1] insidevk/posts/4444.json exists
```
