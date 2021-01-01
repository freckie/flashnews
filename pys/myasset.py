import requests
from bs4 import BeautifulSoup

url = 'https://www.myasset.com/myasset/research/rs_news/news_list.cmd?s_kind_code=00'
req = requests.get(url)
bs = BeautifulSoup(req.text, 'lxml')

wrapper = bs.find('tbody', id='news_tboody')
items = wrapper.find_all('tr')

item_url = 'https://www.myasset.com/myasset/research/rs_news/news_view.cmd?s_kind_code={}&viewdate={}'
for item in items[:3]:
    td = item.find('td', class_='txtL')
    a_tag = td.find('a')
    
    title = a_tag.contents[0].strip()
    date = a_tag['data-seq']

    url2 = item_url.format(a_tag['data-kind'], date.replace(' ', '%20'))
    req2 = requests.get(url2)
    bs2 = BeautifulSoup(req2.text, 'lxml')

    contents_wrapper = bs2.find('dd', class_='view')
    contents = contents_wrapper.text # contents_wrapper 안의 text 노드만 수집

    print(title)
    print(url2)
    print(date)
    print(contents)
