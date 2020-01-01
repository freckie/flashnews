import requests
from bs4 import BeautifulSoup
import chardet

url = 'http://www.mdtoday.co.kr/mdtoday/index.html'

req = requests.get(url)
bs = BeautifulSoup(req.text, 'lxml')

wrapper = bs.find('td', id='MainContent')
items = wrapper.find_all('div', id='box1')

for item in items[:1]:
    a_tag = item.find('a')
    title = bytes(a_tag.find('b').get_text(), 'iso-8859-1').decode('euc-kr')
    href = 'http://m.mdtoday.co.kr' + a_tag['href']
    date = item.find('font').get_text().strip()

    req2 = requests.get(href)
    bs2 = BeautifulSoup(req2.content, 'lxml')
    wrapper2 = bs2.find('div', id='articleBody')
    contents = wrapper2.get_text().strip()

    print(title, href, date, contents)