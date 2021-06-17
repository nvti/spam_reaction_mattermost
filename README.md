# Spam Reaction Mattermost :rofl:

This is a funny project for add large of reaction on a post in mattermost :laughing::laughing::laughing:. Don't use it for bad purpose :scream:

## Usage

```
Usage:  ./react [options] post_link
post_link: Link of post you want to spam reaction =))
Options:
  -email string
        Mattermost login email
  -file string
        File contain mattermost login email and password. Support json and yaml type
  -n int
        Number of reaction you want to add. Set 0 to use max supported emoji (default 20)
  -pass string
        Mattermost login password
Login using certificate file or email+pass
```

For login, you can pass your email and password throw command line argument or you can use a certificate file (Save time for multiple use :rofl:). There are 2 type of supported certificate: json and yaml. You can view content of example file in [sample.json](sample.json) and [sample.yaml](sample.yaml).

About post link, you can get it by click on `More actions` button on a post and click `Copy link`. The link of the post will be copied on clipboard.

![get link](get_link.png)

Also, you can set the number of reaction to use by passing -n argument

Example:

```
./react -f cert.json -n 30 http://localhost/bod/pl/jkbzhub6sbgj7dwkxju9tikjye
```

---

Hope you enjoy this project and have funny moment with your friend :kissing_heart:
