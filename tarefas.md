[ x ] - Criar os topico de controle = Usuários
[] - Criar os topico de controle = Grupos

------

[ x ] - Receber do cli um id

[ x ] - criar um topico novo de controle a partir do id
ID_Control (caso exista se isncreve ao topico)

[ x ] - Mandar uma mensagem para o topico usuários mandando o 
id, status e nomeControl

[ x ] - se inscrever no topico usuários
[ ] - se inscrever no topicos grupos

------


--------

[ ] - enviar solicitação de criação de novo grupo
(olhar na memoria se o grupo existe, caso exista, mandar solicitação ao control_dono_do_grupo ) caso não criar novo grupp

[ ] criar grupo é mandar uma mensagem para o topico de controle de grupos.

--------

--------
[ ] Enviar solicitação para conversar com outro usuário

(isso é enviar uma solicitação ao control do usuario com uma mensagem de "iniciar chat")

[] Receber mensagens de solicitação de novas conversar no seu topico de controle e aceitar ou não
  se aceitar mandar mensagem para o topico do solicitante
  se aceitar se inscrever no topico da novo x_y_timestamp

[ ] Receber mensagens de "aceite de conversação" e se inscrever no grupo

-------------

GRUPOS

-------

[ ] Usuário receber a mensagem de criação de um grupo novo e armazenar a informação em algum lugar

[ ] Usuário solicitar entrar em grupo novo
  - Pegar o dono do grupo
  - Mandar mensagem para o topico de controle do dono do grupo

[ ] Usuário dono de grupo pode aceitar ou não um usuário novo
  - Caso aceite
    [ ] Mandar uma mensagem para o topico de controle do usuário com o tópico do grupo para ele se inscrever
    [ ] Tópico de grupos com o novo usuário do grupo
[ ] Usuário tem que ouvir o aceite de solicitação de entrada no grupo
  - Se inscreve no topico do grupo


--- CLI
thread de cli
[ ] - listar mensagens 
[ ] - listar grupos
	fmt.Println("1 - Listar Usuários")
		fmt.Println("2 - Criar novo grupo")
		fmt.Println("3 - Listar grupos")
		fmt.Println("4 - Nova conversa")
		fmt.Println("5 - Listagem do histórico de solicitação recebidas")
		fmt.Println("6 - Listagem das confirmações de aceitação da solicitação de batepapo")

Abrir um bate papo
Mandar mensagem nova
Ver mensagens de uma conversa
