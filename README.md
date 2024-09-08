# govgr




Testing
=======

You should create a .env file at the root of the project:

gsisUserUsername="myUserName"
gsisUserPassword="#########"


go test -v ./authgovgr

Run all tests in a file: go test -v ./authgovgr/authgovgr_test.go
Run all tests in a package: go test -v ./authgovgr
Run a specific test: go test -v ./authgovgr -run TestLogin
Run all tests in all packages: go test -v ./...
Run tests with coverage: go test -v -cover ./authgovgr

https://www.gov.gr/ipiresies/polites-kai-kathemerinoteta/psephiaka-eggrapha-gov-gr/ekdose-upeuthunes-deloses


https://dilosi.services.gov.gr/login?template=YPDIL&next=/templates/YPDIL/create&

    Citizen: {
      firstname: '',
      surname: '',
      father_fullname: '',
      mother_fullname: '',
      birth_date: '',
      birth_place: '',
      adt: '',
      afm: '',
      residence: '',
      street: '',
      street_number: '',
      tk: '',
      tel: '',
      email: ''
    },

TODO
======

==> Rename packages to remove ambiguity when imported from the end user.
==> Use a randomly choosen up to date User Agent?
==> How to make function signatures like the destructuring style in nodejs i.s {first_value,second_value} so that order of arguments doesnt matter?
==> Is it better to create a class GovGrUser so to avoid passing the userName and userPassword all the time?
==> Automate authentication testing to send an email when something breaks at gov.gr api (i.e. subdomain change).
==> When something breaks make it clear which http call had the problem.
==> Add option to operate behind system proxy.