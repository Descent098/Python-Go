import gc
import time
import requests
import psutil
from urllib.parse import urlparse
from scraping import Site
from concurrent.futures import ThreadPoolExecutor

def get_peak_cpu_usage(func, *args, **kwargs):
    cpu_usages = []
    memory = []
    start_time = time.time()
    
    # Start monitoring CPU usage in a loop
    while (time.time() - start_time) < 1:
      cpu_usages.append(psutil.cpu_percent(interval=0.1)) #Check CPU every 0.1 second
      memory.append(psutil.virtual_memory().percent)
      
    result = func(*args, **kwargs) # Execute the function
    
    # Continue monitoring CPU usage after function execution for 1 second
    start_time = time.time()
    while (time.time() - start_time) < 1:
      cpu_usages.append(psutil.cpu_percent(interval=0.1))
      memory.append(psutil.virtual_memory().percent)
    return max(cpu_usages), max(memory)

def get_data_from_url(url:str) -> Site:
    # print(f"Parsing {url}")
    try:
        temp = urlparse(url)
        resp = requests.get(url, timeout=5)
    except Exception as e:
        print(f"Could not parse {url} due to {e}")
        return
    return Site(
        url, 
        temp.hostname,
        resp.headers.get("server", ""),
        temp.scheme,
        resp.headers.get("Content-Type", "text/html"),
        resp.text,
        temp.port
    )

def pure_python_scraping(urls:list[str]) -> list[Site]:
    with ThreadPoolExecutor(max_workers = 10) as executor:
        temp = executor.map(get_data_from_url, urls)
        result = list(temp)
    return result

def fmt_message(message:str):
    return f"{message:-^15}"

def simple_benchmarking(urls:list[str]):
    print("\x1b[33m")
    try:
        print(fmt_message('Python Starting'))
        t1 = time.time()
        cpu_1, mem_1 = get_peak_cpu_usage(pure_python_scraping, urls)
        t2 = time.time()
        print(fmt_message('Python Ending'))
        
        # Clean up memory so the measurements aren't that tainted
        gc.collect()
        time.sleep(.75) # 750ms should be more than enough time for the gc to finish
        
        print("\x1b[36m")
        print(fmt_message('Go Starting'))
        t3 = time.time()
        cpu_2, mem_2 = get_peak_cpu_usage(Site.from_urls, urls)
        t4 = time.time()
        print(fmt_message('Go ending'))
    finally:
        print("\x1b[37m") # reset to white
    
    time.sleep(2)
    
    print(f"\n{fmt_message('Results')}")    
    print(f"pure_python_scraping took {t2-t1} seconds for {len(urls)} sites with:\n\tMax CPU usage of {cpu_1}\n\tMax RAM usage of: {mem_1} %")
    print(f"Site.from_urls took {t4-t3} seconds for {len(urls)} sites with:\n\tMax CPU usage of {cpu_2}\n\tMax RAM usage of: {mem_2} %")
    total_memory = psutil.virtual_memory().total
    
    difference = max(mem_1, mem_2) - min(mem_1, mem_2)
    difference_multiplier = difference/100
    print(f"Memory difference is {difference}% of {total_memory}\n\t{difference_multiplier*total_memory}B\n\t{(difference_multiplier*total_memory)//1024}KB\n\t{((difference_multiplier*total_memory)//1024)//1024}MB")

if __name__ == "__main__":
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
        "https://www.dropbox.com",
        "https://www.spotify.com",
        "https://www.tumblr.com",
        "https://www.quora.com",
        "https://www.stackoverflow.com",
        "https://www.medium.com",
        "https://www.bing.com",
        "https://www.paypal.com",
        "https://www.ebay.com",
        "https://www.pinterest.com",
        "https://www.tiktok.com",
        "https://www.cnn.com",
        "https://www.bbc.com",
        "https://www.nytimes.com",
        "https://www.theguardian.com",
        "https://www.washingtonpost.com",
        "https://www.forbes.com",
        "https://www.bloomberg.com",
        "https://www.airbnb.com",
        "https://www.udemy.com",
        "https://www.coursera.org",
        "https://www.khanacademy.org",
        "https://www.github.com",
        "https://www.gitlab.com",
        "https://www.codepen.io",
        "https://www.heroku.com",
        "https://www.digitalocean.com",
        "https://www.slack.com",
        "https://www.zoom.us",
        "https://www.skype.com",
        "https://www.trello.com",
        "https://www.notion.so",
        "https://www.canva.com",
        "https://www.wix.com",
        "https://www.shopify.com",
        "https://www.mozilla.org",
        "https://www.icann.org",
        "https://www.cloudflare.com",
        "https://www.openai.com",
        "https://www.deepmind.com",
        "https://www.ibm.com",
        "https://www.oracle.com",
        "https://www.sap.com",
        "https://www.adobe.com",
        "https://www.salesforce.com",
        "https://www.zendesk.com",
        "https://www.asana.com",
        "https://www.bitbucket.org",
        "https://www.bitly.com",
        "https://www.hubspot.com",
        "https://www.mailchimp.com",
        "https://www.figma.com",
        "https://www.behance.net",
        "https://www.dribbble.com",
        "https://www.envato.com",
        "https://www.codeacademy.com",
        "https://www.pluralsight.com",
        "https://www.edx.org",
        "https://www.futurelearn.com",
        "https://www.teachable.com",
        "https://www.skillshare.com",
        "https://www.lynda.com",
        "https://www.x.com",
        "https://www.aliexpress.com",
        "https://www.flipkart.com",
        "https://www.target.com",
        "https://www.homedepot.com",
        "https://www.walmart.com",
        "https://www.bestbuy.com",
        "https://www.nike.com",
        "https://www.adidas.com",
        "https://www.samsung.com",
        "https://www.huawei.com",
        "https://www.sony.com",
        "https://www.lenovo.com",
        "https://www.dell.com",
        "https://www.hp.com",
        "https://www.intel.com",
        "https://www.amd.com",
        "https://www.nvidia.com",
        "https://www.tesla.com",
        "https://www.ford.com",
        "https://www.gm.com",
        "https://www.toyota.com",
        "https://www.honda.com",
        "https://www.bmw.com",
        "https://www.mercedes-benz.com"
    ]
    
    additional_urls = [
    "https://www.ycombinator.com",
    "https://www.producthunt.com",
    "https://www.crunchbase.com",
    "https://www.techcrunch.com",
    "https://www.engadget.com",
    "https://www.theverge.com",
    "https://www.wired.com",
    "https://www.zdnet.com",
    "https://www.cnet.com",
    "https://www.lifehacker.com",
    "https://www.makeuseof.com",
    "https://www.arstechnica.com",
    "https://www.tomshardware.com",
    "https://www.howtogeek.com",
    "https://www.sciencedaily.com",
    "https://www.nature.com",
    "https://www.sciencemag.org",
    "https://www.popsci.com",
    "https://www.space.com",
    "https://www.nasa.gov",
    "https://www.noaa.gov",
    "https://www.who.int",
    "https://www.cdc.gov",
    "https://www.nih.gov",
    "https://www.whitehouse.gov",
    "https://www.congress.gov",
    "https://www.supremecourt.gov",
    "https://www.un.org",
    "https://www.worldbank.org",
    "https://www.imf.org",
    "https://www.oecd.org",
    "https://www.weforum.org",
    "https://www.undp.org",
    "https://www.unesco.org",
    "https://www.ted.com",
    "https://www.brainpickings.org",
    "https://www.goodreads.com",
    "https://www.bookbub.com",
    "https://www.librarything.com",
    "https://www.archlinux.org",
    "https://www.ubuntu.com",
    "https://www.debian.org",
    "https://www.fedoraproject.org",
    "https://www.linuxmint.com",
    "https://www.kali.org",
    "https://www.gentoo.org",
    "https://www.apache.org",
    "https://www.nginx.com",
    "https://www.mysql.com",
    "https://www.postgresql.org",
    "https://www.mongodb.com",
    "https://www.redis.io",
    "https://www.sqlite.org",
    "https://www.rabbitmq.com",
    "https://www.docker.com",
    "https://www.kubernetes.io",
    "https://www.jenkins.io",
    "https://www.travis-ci.com",
    "https://www.circleci.com",
    "https://www.netlify.com",
    "https://www.vercel.com",
    "https://www.render.com",
    "https://www.herokuapp.com",
    "https://www.supabase.com",
    "https://www.prisma.io",
    "https://www.grafana.com",
    "https://www.prometheus.io",
    "https://www.elastic.co",
    "https://www.datadoghq.com",
    "https://www.splunk.com",
    "https://www.cloudflarestatus.com",
    "https://www.cloudflarestatus.com",
    "https://www.wolframalpha.com",
    "https://www.desmos.com",
    "https://www.geogebra.org",
    "https://www.overleaf.com",
    "https://www.latex-project.org",
    "https://www.arxiv.org",
    "https://www.researchgate.net",
    "https://www.jstor.org",
    "https://www.acm.org",
    "https://www.ieee.org",
    "https://www.springer.com",
    "https://www.elsevier.com",
    "https://www.scopus.com",
    "https://www.duolingo.com",
    "https://www.memrise.com",
    "https://www.busuu.com",
    "https://www.babbel.com",
    "https://www.hellotalk.com",
    "https://www.italki.com",
    "https://www.lingq.com",
    "https://www.openstreetmap.org",
    "https://www.mapbox.com",
    "https://www.here.com",
    "https://www.uber.com",
    "https://www.lyft.com",
    "https://www.doordash.com",
    "https://www.ubereats.com",
    "https://www.grubhub.com"
    ]
    
    simple_benchmarking([*urls, *additional_urls])
    
    

