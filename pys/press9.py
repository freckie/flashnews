# paxety랑 같음

import requests
from bs4 import BeautifulSoup

url = 'http://www.press9.kr/news/articleList.html?page=1&sc_section_code=S1N12&sc_order_by=E'

req = requests.get(url)
bs = BeautifulSoup(req.text, 'lxml')

wrapper = bs.find('section', class_='article-list-content')
items = wrapper.find_all('div', class_='table-row')

for item in items[0:10]:
    a_tag = item.find('a', class_='links')
    href = 'http://www.press9.kr' + a_tag['href'].strip()
    title = a_tag.get_text().strip()
    date = item.find("div", class_='list-dated').get_text().split(' | ')[1].strip()

    req2 = requests.get(href)
    bs2 = BeautifulSoup(req2.text, 'lxml')
    
    wrapper2 = bs2.find('div', id='article-view-content-div')
    contents = ''.join([it.get_text().strip() for it in wrapper2.find_all('p')])

    print('=================')
    print(title)
    print(date)
    print(href)
    print(contents)