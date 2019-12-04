import requests
from bs4 import BeautifulSoup

url = 'https://www.mk.co.kr/news/all/'

req = requests.get(url)
bs = BeautifulSoup(req.content, 'lxml')

wrapper = bs.find('div', class_='list_area')
items = wrapper.find_all('dl', class_='article_list')

for item in items[:1]:
    dt = item.find('dt', class_='tit')

    a_tag = dt.find('a')
    title = a_tag.get_text().strip()
    href = a_tag['href']
    date = item.find('span', class_='date').get_text().strip()

    req2 = requests.get(href)
    bs2 = BeautifulSoup(req2.content, 'lxml')
    wrapper2 = bs2.find('div', class_='art_txt')
    removes = wrapper2.find_all('div')
    removes.extend(wrapper2.find_all('script'))

    contents = wrapper2.get_text().strip()
    for remove in removes:
        contents = contents.replace(remove.get_text(), '')

    print(title, href, date, contents)