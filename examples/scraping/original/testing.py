from scraping import Site

print(Site.from_str("https://kieranwood.ca"))

print(Site.from_urls(["https://google.ca", "https://cloudflare.ca"]))

urls = [
        "https://www.google.com",
        "https://www.facebook.com",
        "https://www.youtube.com",
        "https://www.twitter.com",
        "https://www.instagram.com",
        "https://www.linkedin.com",
        "https://www.wikipedia.org",
        "https://www.reddit.com",
        "https://www.amazon.com",
        "https://www.netflix.com",
        "https://www.apple.com",
        "https://www.microsoft.com",
    ]

r:list[Site] = Site.from_urls(urls)

print(len(r))
for site in r:
    print(site.url)
