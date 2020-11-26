# doctorsnews랑 같음

import requests
from bs4 import BeautifulSoup

url = 'http://marketinsight.hankyung.com/apps.free/free.news.list?category=IB_FREE'

req = requests.get(url)
bs = BeautifulSoup(req.text, 'lxml')

wrapper = bs.find('ul', id='newsList')
items = wrapper.find_all('li', recursive=False)
print(items)


req2 = requests.get('http://marketinsight.hankyung.com/apps.free/free.news.view?aid=202011266458u&category=IB_FREE')
bs2 = BeautifulSoup(req2.text, 'lxml')
print(bs2)

exit(1)

for item in items[0:3]:
    a_tag = item.find('a', class_='dyn std')
    href = 'https://dealsite.co.kr' + a_tag['href'].strip()
    title = a_tag.get_text().strip()
    date = item.find("div", class_='pubdate').get_text().strip()
    contents = item.find('span', class_='sneakpeek').get_text().strip()

    print('=================')
    print(title)
    print(date)
    print(href)
    print(contents)
