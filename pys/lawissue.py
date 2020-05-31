import requests
from bs4 import BeautifulSoup

url = 'https://www.lawissue.co.kr/list.php?ct=g0000&ssk=&nmd=2'
req = requests.get(url)
bs = BeautifulSoup(req.text, 'lxml')

wrapper = bs.find('ul', class_='lst_type_01')
items = wrapper.find_all('li')

for item in items[:15]:
    a_tag = item.find('a', class_='tit')

    title = a_tag.get_text().strip()
    url = '' + a_tag['href']
    date = item.find('span', class_='date').get_text().split('|')[1]

    req2 = requests.get(url)
    bs2 = BeautifulSoup(req2.text, 'lxml')
    wrapper2 = bs2.find('div', id='CmAdContent')

    contents = ''.join([it.strip() for it in wrapper2.find_all(text=True, recursive=False)])

    print(title)
    print(url)
    print(date)
    print(contents)