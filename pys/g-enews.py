# thebell 거의 일치

import requests
from bs4 import BeautifulSoup

url = 'http://m.g-enews.com/issuelist.php?ud=2017011301560109486&ct=g000000'
req = requests.get(url)
bs = BeautifulSoup(req.text, 'lxml')

wrapper = bs.find('ul', id='liAppendID')
items = wrapper.find_all('li')

for item in items[:1]:
    a_tag = item.find('a')
    href = '' + a_tag['href']
    title = item.find('span', class_='elip2').get_text().strip()
    date = item.find('span', class_='r2').get_text().strip()

    req2 = requests.get(href)
    bs2 = BeautifulSoup(req2.text, 'lxml')

    wrapper2 = bs2.find('div', id='articleBody')
    contents = ''.join([it.strip() for it in wrapper2.find_all(text=True, recursive=False)])

    print(title, href, date, contents)