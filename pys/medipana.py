# ddaily 거의 일치

import requests
from bs4 import BeautifulSoup

url = 'https://www.medipana.com/news/news_list_new.asp?MainKind=A&NewsKind=103&vCount=15&vKind=1'
req = requests.get(url)
bs = BeautifulSoup(req.text, 'lxml')

wrapper = bs.find('div', class_='totalNews')
items = wrapper.find_all('li')

for item in items[:1]:
    a_tag = item.find('a')
    href = 'https://www.medipana.com/news/' + a_tag['href']
    title = item.find('span', class_='tit').get_text().strip()
    date = item.find('span', class_='infor').get_text().split(' | ')[1].strip()

    req2 = requests.get(href)
    bs2 = BeautifulSoup(req2.text, 'lxml')

    wrapper2 = bs2.find('div', class_='newsCon')
    contents = '\n'.join([it.get_text().strip() for it in wrapper2.find_all('div')])

    print(title, href, date, contents)