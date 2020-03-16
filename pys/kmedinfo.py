import requests
from bs4 import BeautifulSoup

url = 'http://m.kmedinfo.co.kr'

req = requests.get(url)
bs = BeautifulSoup(req.text, 'lxml')

wrapper = bs.find('div', id="set_A1Container").find('ul')
items = wrapper.find_all('li')

result = list()
for item in items:
    temp = {
        'title': '',
        'link': '',
        'date': '',
        'content': ''
    }
    a_tag = item.find('a')

    temp['title'] = a_tag.get_text().strip()
    temp['link'] = 'http://m.kmedinfo.co.kr' + a_tag['href']

    # contents
    req2 = requests.get(temp['link'])
    bs2 = BeautifulSoup(req2.text, 'lxml')
    wrapper2 = bs2.find('div', id='articleBody').find('div', class_='body')
    temp['date'] = bs2.find('p', class_='date').get_text().replace('기사승인', '').replace('\xa0', ' ').strip()
    temp['content'] = ' '.join([it.get_text().strip() for it in wrapper2.find_all('p')])

    result.append(temp)

print(result)