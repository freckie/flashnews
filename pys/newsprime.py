# doctorsnews와 일치

import requests
from bs4 import BeautifulSoup

url = 'http://m.newsprime.co.kr/section_list.html?sec_no=56&menu=index'
req = requests.get(url)
bs = BeautifulSoup(req.content, 'lxml')

wrapper = bs.find('div', class_='box01_0610_section')
items = wrapper.find_all('div', class_='article_box_sl_section')

for item in items[:15]:
    a_tag = item.find('div', class_='title_text').find('a')
    href = 'http://m.newsprime.co.kr/' + a_tag['href']
    title = a_tag.get_text().strip()

    req2 = requests.get(href)
    bs2 = BeautifulSoup(req2.text, 'lxml')

    date = bs2.find('div', class_='hd').find('p', class_='data').get_text().strip()
    wrapper2 = bs2.find('div', class_='stit2')
    contents = ''
    for it in wrapper2.find_all('div', recursive=False):
        if not it.find('div', class_='imgframe'):
            contents += it.get_text().strip()

    print('==========================')
    print(title)
    print(href)
    print(date)
    print(contents)