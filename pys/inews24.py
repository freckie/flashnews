# yna와 유사

import requests
from bs4 import BeautifulSoup

url = 'http://www.inews24.com/list/inews'
req = requests.get(url)
bs = BeautifulSoup(req.text, 'lxml')

wrapper = bs.find('article', class_='list')
items = wrapper.find_all('li', class_='list')

for item in items[:1]:
    a_tag = item.find('a')
    href = 'http://www.inews24.com' + a_tag['href']
    title = a_tag.get_text().strip()
    date = item.find('time').get_text().strip()

    req2 = requests.get(href)
    bs2 = BeautifulSoup(req2.text, 'lxml')

    wrapper2 = bs2.find('article', id='articleBody')
    contents = ''.join([it.get_text().strip() for it in wrapper2.find_all('p')])

    print(title, href, date, contents)