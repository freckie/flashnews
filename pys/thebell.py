# doctorsnews와 거의 일치

import requests
from bs4 import BeautifulSoup

url = 'http://www.thebell.co.kr/free/content/article.asp?svccode=00'
req = requests.get(url)
bs = BeautifulSoup(req.text, 'lxml')

wrapper = bs.find('div', class_='listBox')
items = wrapper.find_all('li')

for item in items[:1]:
    a_tag = item.find('a')
    href = 'http://www.thebell.co.kr/free/content/' + a_tag['href']
    title = item.find('dt').get_text().strip()
    date = item.find('span', class_='date').get_text().strip()

    req2 = requests.get(href)
    bs2 = BeautifulSoup(req2.text, 'lxml')

    wrapper2 = bs2.find('div', id='article_main')
    contents = ''.join([it.strip() for it in wrapper2.find_all(text=True, recursive=False)])

    print(title, href, date, contents)