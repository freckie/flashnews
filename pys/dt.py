# thebell 거의 일치
# 인코딩 : iso-8859-1로 인코딩 후 euc-kr로 디코딩

import requests
from bs4 import BeautifulSoup

url = 'http://www.dt.co.kr/section.html?section_num=2900'
req = requests.get(url)
bs = BeautifulSoup(req.text, 'lxml')

wrapper = bs.find('div', class_='list_area')
items = wrapper.find_all('dl', class_='article_list')

for item in items[:3]:
    dt = item.find('dt')
    a_tag = dt.find('a')
    href = '' + a_tag['href']
    title = a_tag.get_text().strip().encode('iso-8859-1').decode('euc-kr')
    date = item.find('span', class_='date').get_text().encode('iso-8859-1').decode('euc-kr').replace('입력 ', '').strip()

    req2 = requests.get(href)
    bs2 = BeautifulSoup(req2.text, 'lxml')

    wrapper2 = bs2.find('div', class_='art_txt')
    contents = ''.join([it.strip() for it in wrapper2.find_all(text=True, recursive=False)]).encode('iso-8859-1').decode('euc-kr')

    print(title, href, date, contents)