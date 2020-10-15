package deploy

// https://habr.com/ru/company/ruvds/blog/438982/
// Kubernetes это система для автоматизации развёртывания, масштабирования и управления
// контейнеризированными приложениями (оркестратор контейнеров)

// контейнер это легковесный, автономный, исполняемый пакет содержащий приложение,
// включая все необходимое для его запуска: код, среду исполнения, системные средства и библиотеки, настройки
// при этом все будет одинаково работать как в видовс, так и в линух

// недостаток виартуальных машин в сравнении с контейнерами
// - это ресурсы, т.к. каждая вм требует полноценной операционной системы
// - зависимость от платформы
// - требовательное к ресурсам масштабирование решения основанного на вм

// сильные стороны контейнеров
// - эффективное использование ресурсов
// - независимость от платформы. Контейнер, которые разраб запустит на своем пк будет работать
//                                                                                  где угодно
// - легковесное развертывание за счет слоев образов

// грубо говоря, если ВМ используют каждая свою ОС, то контейнеры в это время используют
// хостовую ОС, "разделяя" ее
// для докера основной файл это Dockerfile
// вначале делают запись о базововом контейнере а потом последовательные инструкции
// на порядок создания контейнера, который будет соответствовать нуждам приложения

// FROM - имя базового контейнар, потом инструкции

// docker login
// если ошибка ... permission denied while trying to connect to the Docker daemon socket at unix
// тогда sudo chmod 666 /var/run/docker.sock

// docker build -f Dockerfile -t {{dockerID}}}}/{{container name}} .
// для отправки в репозиторий
// docker push  {{dockerID}}}}/{{container name}}

// что бы собрать и запустить, любой может
// docker pull {{dockerID}}}}/{{container name}}
// docker run -d -p 80:80 {{dockerID}}}}/{{container name}}
//  80:80  - первая 80 порт хоста, вторая порт контейнера

// docker build требует конекст сборки, т.е. файлы и папки, по которым будет строится контейнер
// но загружать то, что не нужно - это пустая трата времени
// то, что не нужно для сборки указывается в .dockerignore (должен находится там же где докерфайл)

// ENV добавляет перменные окружения
// пример:
// ENV SA_LOGIC_API_URL http://localhost:5000

// EXPOSE указывает, что порт нужно открыть (команда только для документации)
// пример:
// EXPOSE 8080

// docker ps - список запущенных контейнеров
// docker kill {{name}} убить контейнер по имени

// найти ип контейнера
// docker inspect -f \
//  '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' \
//   {{containerID}}

// KUBERNETES позволяет приложениям абстрагироваться от инфраструктуры
// предоставляя АПИ, которому отправляются запросы
// при этом кубернетес старается выполнить используя свои возможности
// чет -типа кубер, дай 4 контейнера, кубер найдет не сильно нагруженные
// ноды, на котрых можно развернуть контейнеры

// API Server мастер-нода на которой развернут кубер называется API Server
// выполнение обращений к этому серверу - единственный способ взаимодействия с кластером (подов)
// если речь идет о остановке и запуске контейнеров, о проверке состояния системы, работе
// с логами и т.д.

// Kubelet это агент, который осуществляет мониторинг контейнеров(подов),находящихся внутри узла
// и который взаимоедйствует с главным узлом

// также кубер способствует стандартизации работы с провайдерами облачных узлов
// т.е. на любом облаке (азур, авс, гугл клоуд платформ) работать с кубер будет одинаково.
// разработчик в декларативном стиле сообщает API Server что ему нужно, а сам кубер
// работает с платформой

// Minikube для управления кластера с одним узлом
// minikube is local Kubernetes, focusing on making it easy to learn and develop for Kubernetes

// install
// curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube_latest_amd64.deb
// sudo dpkg -i minikube_latest_amd64.deb

// PODs
// https://kubernetes.io/docs/concepts/workloads/pods/
// A Pod's contents are always co-located and co-scheduled, and run in a shared context.
// в состав пода могут входить несколько контейнеров, которые используют одну и
// ту же среду выполнения. Но, как правило, один под - один контейнер.
// но если контейнарм нужен доступ к одному и тому-же локальному хранилищу данных
// или между нами налаженно межпроцессорное взаимодействие, любая другая тесная связь,
// это все можно запустить в одном поде
// В подах не обязательно использовать именно контейнеры докер

// Свойства подов
// - у каждого есть свой ИП
// - в поде может быть несколько контейнеров, могут использовать доступные порты (через локалхост)
//   взаимодействие между подами по ип
// - контейнеры в подах совместно используют тома хранилищ данных, ип адрес, номера порыгвтов, пространство имен IPC

// немного о пространстве имен при IPC ( межПРОЦЕССное взаимодействие)
// https://ru.wikipedia.org/wiki/%D0%9F%D1%80%D0%BE%D1%81%D1%82%D1%80%D0%B0%D0%BD%D1%81%D1%82%D0%B2%D0%BE_%D0%B8%D0%BC%D1%91%D0%BD_(Linux)
// пространство имен лежит в основе технологии контейнеров
// это дает процессам запушенным в контейнерах иллюзию, что они имеют собственные ресурсы
// основная цель изоляции процессов состоит в предотвращении вмешательства процессов одного
// контейнера в работу другого а также в работу хостовой машины.

// Контейнеры имеют свои изолированные файловые системы, но могут совместно использовать данные,
// пользуясь ресурсом кубера, который называется volume

// Для создания пода его описывают для каждого приложения в manifest file

/*

apiVersion: v1
kind: Pod                    - вид ресурса кубера
metadata:
 name: sa-frontend           - имя ресурса
spec:                        - объект описывает нужное состояние ресурса
 containers:                 - массив контейнеров
      - image: {{dockerID}}/{{containerID}}        - образ контейнера
      name: sa-frontend                            - уникальное имя контейнера
      ports:                                       - порт, который слушает контейнер (параметр только для документации)
      - containerPort: 80

*/

// создание пода:
// для начала
// minikube start
// kubectl get nodes
// kubectl create -f {{pod_manifest_filename}}.yaml

// если возникают ошибки при запуске, например нечто типа failed to ssh connection
// стоит проверить, как установлен докер
// dpkg -l | grep docker
// если команда ничего не выводит
// и докер при этом установлен через снап, например, нужно его удалить
// и поставить согласно инструкции по ссылке
// https://docs.docker.com/engine/install/ubuntu/#install-using-the-repository
// если не вышло - снести опять все нафиг
// и сделать в точности как здесь
// https://phoenixnap.com/kb/install-minikube-on-ubuntu
// угрохал, бля, часа 4 что бы поставить на убунту