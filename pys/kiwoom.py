import requests
from bs4 import BeautifulSoup

url = 'https://www2.kiwoom.com/nkw.TD6000News.do'
contentURL = 'https://www2.kiwoom.com/nkw.TD6000NewsCont.do'

req = requests.get(url)
bs = BeautifulSoup(req.text, 'lxml')

wrapper = bs.find('table', id='oTable')
items = wrapper.find_all('tr')[1:]

for idx, item in enumerate(items[0:5]):
    a_tag = item.find('td', class_='ldata').find('a')
    title = a_tag.get_text().strip()
    supplier = item.find('div', id='gubn_'+str(idx)).get_text().strip() + \
        item.find('div', id='subg_'+str(idx)).get_text().strip()
    date = item.find('div', id='date_'+str(idx)).get_text().strip() + ' ' + \
        item.find('div', id='time_'+str(idx)).get_text().strip()
    hcode = date.replace('/', '').replace(':', '').replace(' ', '') + item.find('div', id='seqn_'+str(idx)).get_text().strip()
    
    req2 = requests.post(contentURL, data={
        'supplier': supplier,
        'hcode': hcode
    })
    bs2 = BeautifulSoup(req2.text, 'lxml')
    print(bs2)
    contents = bs2.get_text().strip()

    # print('=================')
    # print(item)
    # print(title)
    # print(date)
    # print(hcode)
    # print(supplier)
    # # print(href)
    # print(contents)

    # supplier : 2103
    # hcode : 2020 12 09 13 53 27 13532745
    # title 안 보내도 됨
    # form data로 https://www2.kiwoom.com/nkw.TD6000NewsCont.do POST