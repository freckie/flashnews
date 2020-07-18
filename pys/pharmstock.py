import json
import requests

from urllib import parse
from bs4 import BeautifulSoup

url = 'http://m.pharmstock.co.kr/news/putNewsJson.php'
req = requests.post(url, data={
    'page': 1,
    'sc_order_by': 'E'
})
items = req.json()['news']

for item in items:
    href = 'http://m.pharmstock.co.kr/news/articleView.html?idxno=' + item['idxno']
    title = parse.unquote(item['title']).replace('+', ' ')
    date = item['view_date'] + ' ' + item['view_time']

    req2 = requests.get(href)
    bs2 = BeautifulSoup(req2.text, 'lxml')

    wrapper2 = bs2.find('div', class_='body word_break')
    contents = ''.join([it.get_text().strip() for it in wrapper2.find_all('p', recursive=False)])

    print('=================')
    print(title)
    print(date)
    print(href)
    print(contents)