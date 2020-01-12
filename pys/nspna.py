# 

import requests
from bs4 import BeautifulSoup

url = 'http://www.nspna.com/news/?cid=20'
req = requests.get(url)
bs = BeautifulSoup(req.text, 'lxml')

wrapper = bs.find('div', id='news_list')
items = wrapper.find_all('div', class_='news_panel')

for item in items[:1]:
    a_tag = item.find('a')
    href = 'http://www.nspna.com' + a_tag['href']
    title = item.find('div', class_='subject').get_text().strip()
    date = item.find('div', class_='info').get_text().split(' | ')[1].strip()

    req2 = requests.get(href)
    bs2 = BeautifulSoup(req2.text, 'lxml')

    wrapper2 = bs2.find('div', id='CmAdContent')
    contents = wrapper2.find('p', class_='section_txt').get_text().strip()

    print(title, href, date, contents)