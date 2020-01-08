# doctorsnews와 거의 일치

import requests
from bs4 import BeautifulSoup

url = 'http://www.seoulwire.com/news/articleList.html?view_type=sm'
req = requests.get(url)
bs = BeautifulSoup(req.text, 'lxml')

wrapper = bs.find('section', class_='article-list-content')
items = wrapper.find_all('div', class_='list-block')

for item in items[:1]:
    a_tag = item.find('a')
    href = 'http://www.seoulwire.com' + a_tag['href']
    title = a_tag.get_text().strip()
    date = item.find('div', class_='list-dated').get_text().split(' | ')[2].strip()

    req2 = requests.get(href)
    bs2 = BeautifulSoup(req2.text, 'lxml')

    wrapper2 = bs2.find('div', id='article-view-content-div')
    contents = ''.join([it.get_text().strip() for it in wrapper2.find_all('p')])

    print(title, href, date, contents)