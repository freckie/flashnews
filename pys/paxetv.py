# doctorsnews와 일치

import requests
from bs4 import BeautifulSoup

url = 'http://www.paxetv.com/news/articleList.html?page=1&total=14633&sc_section_code=S1N14'
req = requests.get(url)
bs = BeautifulSoup(req.content, 'lxml')

wrapper = bs.find('section', class_='article-list-content')
items = wrapper.find_all('div', class_='table-row')

for item in items[:15]:
    a_tag = item.find('a', class_='links')
    href = 'http://www.paxetv.com/' + a_tag['href']
    title = a_tag.get_text().strip()
    date = bs.find('div', class_='list-dated').get_text().split(' | ')[1].strip()

    req2 = requests.get(href)
    bs2 = BeautifulSoup(req2.text, 'lxml')

    wrapper2 = bs2.find('div', id='article-view-content-div')
    contents = ''.join([it.get_text().strip() for it in wrapper2.find_all('p')])

    print('==========================')
    print(title)
    print(href)
    print(date)
    print(contents)