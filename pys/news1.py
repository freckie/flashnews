import json
import requests
from bs4 import BeautifulSoup

url = 'http://www.news1.kr/ajax/ajax.php'

params = {
    'cmd': 'categories',
    'op': 'categories_list',
    'slimit': 1,
    'elimit': 15,
    'orderby': 'pubdate_tsm',
    'sort': 'DESC',
    'categories_sec': 'parent',
    'upper_categories_id': '13',
    'categories_id': '13'
}
req = requests.post(url, data=params)
json_data = req.json()

for item in json_data['data']:
    title = item['title']
    href = 'http://www.news1.kr/articles/?' + item['id']
    date = item['6']
    bs = BeautifulSoup(item['3'], 'lxml')
    contents = bs.get_text().replace(bs.find('table').get_text(), '')

    print(title, href, date, contents)