import requests
from bs4 import BeautifulSoup

url = 'http://biz.heraldcorp.com/list.php?ct=010106000000&ctm=19'
req = requests.get(url)
bs = BeautifulSoup(req.text, 'lxml')

items = []

first = bs.find('div', class_='list_head')
first_a = first.find('a')
first_title = first.find('div', class_='list_head_title').get_text().strip()
first_date = first.find('div', class_='list_date').get_text().strip()
first_url = 'http://biz.heraldcorp.com/' + first_a['href']

wrapper = bs.find('div', class_='list').find('ul')
items = wrapper.find_all('li')

for item in items[:3]:
    a_tag = item.find('a')
    url2 = 'http://biz.heraldcorp.com/' + a_tag['href']
    
    title = a_tag.find('div', class_='list_title').get_text().strip()
    date = item.find('div', class_='l_date').get_text().strip()

    req2 = requests.get(url2)
    bs2 = BeautifulSoup(req2.text, 'lxml')

    contents_wrapper = bs2.find('div', id='articleText')
    contents = ' '.join([it.get_text().strip() for it in contents_wrapper.find_all('p')])

    print(title)
    print(url2)
    print(date)
    print(contents)
