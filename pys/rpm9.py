# doctorsnews와 일치

import requests
from bs4 import BeautifulSoup

url = 'http://www.rpm9.com/news/section.html?id1=6'
req = requests.get(url)
bs = BeautifulSoup(req.content, 'lxml')

wrapper = bs.find('ul', class_='sub_newslist')
items = wrapper.find_all('li')

for item in items[:10]:
    if item.has_attr('class') and item['class'][0] == 'ad_text660':
        continue

    a_tag = item.find('a', class_='newstit')
    href = 'http://www.rpm9.com' + a_tag['href']
    title = a_tag.get_text().strip()

    req2 = requests.get(href)
    bs2 = BeautifulSoup(req2.text, 'lxml')

    date = bs2.find('span', class_='date').get_text().strip().replace('발행일 : ', '')
    wrapper2 = bs2.find('div', id='articleBody')
    contents = ''.join([it.get_text().strip() for it in wrapper2.find_all('p')])

    print('==========================')
    print(title)
    print(href)
    print(date)
    print(contents)