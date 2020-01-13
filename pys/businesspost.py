# doctorsnews와 거의 일치

import requests
from bs4 import BeautifulSoup

url = 'http://m.businesspost.co.kr/BP?command=mobile_list&sc_cate=ALL'
req = requests.get(url)
bs = BeautifulSoup(req.text, 'lxml')

wrapper = bs.find('ul', class_='list_thumbandstat')
items = wrapper.find_all('li')

for item in items[:5]:
    a_tag = item.find('a')
    href = 'http://m.businesspost.co.kr' + a_tag['href']
    title = item.find('strong').get_text().strip()

    req2 = requests.get(href)
    bs2 = BeautifulSoup(req2.text, 'lxml')

    date = bs2.find('div', class_='info') # utils.GetDateString 사용해서 추출
    wrapper2 = bs2.find('div', class_='post-contents')
    contents = ''.join([it.strip() for it in wrapper2.find_all(text=True, recursive=False)])

    print(title, href, date, contents)