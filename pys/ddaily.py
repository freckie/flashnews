# doctorsnews와 거의 일치

import requests
from bs4 import BeautifulSoup

url = 'http://www.ddaily.co.kr/news/article_list_all/'
req = requests.get(url)
bs = BeautifulSoup(req.text, 'lxml')

wrapper = bs.find('div', class_='m01_ara')
items = wrapper.find_all('dl')

for item in items[:3]:
    a_tag = item.find('a')
    href = 'http://www.ddaily.co.kr' + a_tag['href']
    title = item.find('dt').get_text().strip()
    date = ''.join(item.find('span').get_text().strip().split(' ')[1:])

    req2 = requests.get(href)
    bs2 = BeautifulSoup(req2.text, 'lxml')

    wrapper2 = bs2.find('div', id='news_body_area')
    contents = '\n'.join([it.get_text().strip() for it in wrapper2.find_all('div')])

    print(title, href, date, contents)