# doctorsnews와 거의 일치

import requests
from bs4 import BeautifulSoup

url = 'http://www.newsway.co.kr/news/lists'
req = requests.get(url)
bs = BeautifulSoup(req.text, 'lxml')

wrapper = bs.find('div', class_='totalList')
items = wrapper.find_all('li')

for item in items[:1]:
    div = item.find('div', class_='ritext')
    a_tag = div.find('a')
    href = 'http://www.newsway.co.kr' + a_tag['href']
    title = a_tag.find('strong').get_text().strip()
    date = div.find('span').get_text().strip()

    req2 = requests.get(href)
    bs2 = BeautifulSoup(req2.text, 'lxml')

    wrapper2 = bs2.find('div', id='articleText')
    contents = ''.join([it.strip() for it in wrapper2.find_all(text=True, recursive=False)])

    print(title, href, date, contents)