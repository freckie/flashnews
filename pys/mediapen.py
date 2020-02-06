# paxetv 일치

import re
import requests
from bs4 import BeautifulSoup

url = 'http://m.mediapen.com/news/lists/?menu=2&cate_cd1='
req = requests.get(url)
bs = BeautifulSoup(req.content, 'lxml')

wrapper = bs.find('ul', class_='main-list')
items = wrapper.find_all('li')

for item in items[:15]:
    if item.has_attr('class') and item['class'][0] == 'ad_text660':
        continue

    a_tag = item.find('a')
    href = 'http://m.mediapen.com' + a_tag['href']
    title = a_tag.get_text().strip()
    title = re.sub('[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}', '', title).replace('\n', '')

    req2 = requests.get(href)
    bs2 = BeautifulSoup(req2.text, 'lxml')

    date = bs2.find('a', class_='news_author').find_all(text=True)[1].strip()

    wrapper2 = bs2.find('div', class_='news_contents')
    contents = ''.join([it.get_text().strip() for it in wrapper2.find_all('p')])

    print('==========================')
    print(title)
    print(href)
    print(date)
    print(contents)