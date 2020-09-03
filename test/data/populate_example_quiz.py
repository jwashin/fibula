import sys

#sys.path.extend(['/usr/lib/python2.7/site-packages'])

#import logging
##logger = logging.getLogger("lovely.jsonrpc.dispatcher")
#logger.addHandler(logging.StreamHandler(sys.stderr))
#logger.setLevel(logging.DEBUG)

#for item in sorted(sys.path):
#    print item
#from ZODB.FileStorage import FileStorage
#from ZODB.DB import DB
#import jsonrpclib
from operator import attrgetter
import json
import os
from jsonrpc.proxy import JSONRPCProxy
client = JSONRPCProxy.from_url('http://localhost:8080/rpc/login')
#client = JSONRPCProxy.from_url('https://gae_project.appspot.com/rpc/login')
#client = JSONRPCProxy.from_url('http://localhost:8080/xadmin')
#f = client.reverse('hello world!')
#print f
authenticated = client.login('username', 'password')

if authenticated:
    print("authentication OK")
else:
    print("authentication NOT OK")

client._path = 'xadmin'
event_key = client.create_fbla_event(title='Example exam', organization='', proctoring=False, registering=False)

for item in client.list_all_events():
    print ("{}".format(item))

exams_location = '~/exams_test/'

contest_info = json.load(open('test_quiz_data.json', 'r'))

client.event_set_exam_metadata(event_key,contest_info)

for item in contest_info['exams']:
    exam_id = item['id']
    print ("adding {}".format(exam_id))
    with open(os.path.join(exams_location, "{}.json".format(exam_id))) as jsonfile:
        data = json.load(jsonfile)
        exam_key = client.event_add_exam(event_key,data)
        client.exam_set_params(exam_key=exam_key,time=5,need_registration=False)

import time
time.sleep(2)

test_exam_name = client.get_exam_name(event_key,'test_quiz')
assert test_exam_name == 'Example online quiz'
print ("Exam name OK: {}".format(test_exam_name))

#client.event_set_registration_start(event_key, 2012, 3, 21)
#client.event_set_registration_end(event_key, 2012, 3, 27)
#client.event_set_competition_start(event_key, 2012, 3, 28)
#client.event_set_competition_end(event_key, 2012, 4, 11)

event_modes = ['closed', 'development', 'open for registration',
               'examinations in progress', 'results']

client.event_set_mode(event_key, event_modes[3])

school_key = client.create_school(event_key, "Example School", "Example Region")

#school_keys = {}
#school_key = ''
#for region in root.regions[0:1]:
#    #schools = [school for school in region.schools if len(school.students)]
#    schools = sorted(
#        [school for school in region.schools if len(school.students)],
#        key=attrgetter('name'))
#    schools = schools[1:4]
#    for school in schools:
#        school_key = client.create_school(event_key, school.name, region.name)
#        school_keys[(region.name,school.name)] = school_key
#        print school_key, region.name, school.name
#        t = school.getProctors()
#        for proctor in t:
#            print(proctor)
#        for student in school.students:
#            sid = client.school_add_student(school_key, student.last, student.first, student.grade_level)
#            for testtaker in school.test_takers:
#                exam_id = testtaker.exam_name
#                for tt_student in testtaker.students:
#                    if (tt_student.last == student.last and
#                        tt_student.first == student.first and
#                        tt_student.grade_level == student.grade_level):
#                        ok = client.school_register_student(school_key, sid, exam_id)
#
#for item in client.school_get_credentials(school_key):
#    print (item)
#
#for item in client.event_list_schools(event_key):
#    print(item)
#
#
#connection.close()
