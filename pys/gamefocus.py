import requests
from bs4 import BeautifulSoup

url = 'http://www.gamefocus.co.kr/html_file.php?file=normal_all_news.html'
req = requests.get(url)
bs = BeautifulSoup(req.content, 'lxml')

wrapper = bs.find('div', class_='f_l').find('table')
items = wrapper.find_all('tr', recursive=False)

for idx, item in enumerate(items[:15]):
    a_tag = item.find('a')
    href = 'http://www.gamefocus.co.kr/' + a_tag['href']

    req2 = requests.get(href)
    bs2 = BeautifulSoup(req2.text, 'lxml')

    title = bs2.find('div', class_='detail_view').find('h2').get_text().strip()
    date = bs2.find('span', class_='font_11').get_text().strip()
    wrapper2 = bs2.find('div', id='ct')
    contents = ''.join([it.get_text().strip() for it in wrapper2.find_all('p')])

    print('==========================')
    print(idx)
    print(title)
    print(href)
    print(date)
    
    print(contents)