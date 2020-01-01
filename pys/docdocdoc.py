import json
import requests
from bs4 import BeautifulSoup
from urllib import parse

url = 'http://m.docdocdoc.co.kr/news/putNewsJson.php'

params = {
    'page': 1,
    'list_per_page': 10
}
req = requests.post(url, data=params)
json_data = req.json()

for item in json_data['news']:
    title = parse.unquote(item['title']).replace('+', ' ')
    href = 'http://m.docdocdoc.co.kr/news/articleView.html?idxno=' + item['idxno']
    date = item['view_date']

    req = requests.get(href)
    bs = BeautifulSoup(req.text, 'lxml')
    wrapper = bs.find('div', id='articleBody')

    p_tags = wrapper.find_all('p')
    contents = ''
    for p in p_tags:
        contents = contents + p.get_text()

    print(title, href, date, contents)