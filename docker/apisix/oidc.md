# openid-connect

## Descrição
O Plugin `openid-connect` suporta a integração com provedores de identidade OpenID Connect (OIDC), como Keycloak, Auth0, Microsoft Entra ID, Google, Okta, entre outros. Ele permite que o APISIX autentique clientes e obtenha suas informações do provedor de identidade antes de permitir ou negar seu acesso aos recursos protegidos *upstream*.

## Atributos

| Nome | Tipo | Obrigatório | Padrão | Valores Válidos | Descrição |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `client_id` | string | Verdadeiro | | | ID do cliente OAuth. |
| `client_secret` | string | Verdadeiro | | | Segredo do cliente OAuth. |
| `discovery` | string | Verdadeiro | | | URL para o documento de descoberta *well-known* do provedor OpenID, que contém uma lista de endpoints da API do Provedor OIDC (OP). O Plugin pode utilizar diretamente os endpoints deste documento. Você também pode configurar estes endpoints individualmente, o que terá precedência sobre os endpoints fornecidos no documento de descoberta. |
| `scope` | string | Falso | `openid` | | Escopo OIDC que corresponde às informações que devem ser retornadas sobre o usuário autenticado, também conhecidas como *claims*. É usado para autorizar usuários com a permissão adequada. O valor padrão é `openid`, o escopo necessário para o OIDC retornar uma *claim* `sub` que identifica unicamente o usuário autenticado. Escopos adicionais podem ser anexados e delimitados por espaços, como `openid email profile`. |
| `required_scopes` | array[string] | Falso | | | Escopos obrigatórios a estarem presentes no *access token*. Usado em conjunto com o endpoint de introspecção quando `bearer_only` é `true`. Se algum escopo obrigatório estiver faltando, o Plugin rejeitará a solicitação com um erro `403 forbidden`. |
| `realm` | string | Falso | `apisix` | | *Realm* no cabeçalho de resposta `WWW-Authenticate` que acompanha uma solicitação `401 unauthorized` devido a um *bearer token* inválido. |
| `bearer_only` | boolean | Falso | `false` | | Se `true`, exige estritamente um *access token* do tipo *bearer* nas solicitações para autenticação. |
| `logout_path` | string | Falso | `/logout` | | Caminho para ativar o logout. |
| `post_logout_redirect_uri` | string | Falso | | | URL para redirecionar os usuários após o `logout_path` receber uma solicitação para sair. |
| `redirect_uri` | string | Falso | | | URI para redirecionar após a autenticação com o provedor OpenID. Observe que a URI de redirecionamento não deve ser a mesma que a URI da solicitação, mas um subcaminho da URI da solicitação. Por exemplo, se o `uri` da Rota for `/api/v1/*`, `redirect_uri` pode ser configurado como `/api/v1/redirect`. Se `redirect_uri` não estiver configurado, o APISIX acrescentará `/.apisix/redirect` à URI da solicitação para determinar o valor de `redirect_uri`. |
| `timeout` | integer | Falso | `3` | `[1,...]` | Tempo de *timeout* da solicitação em segundos. |
| `ssl_verify` | boolean | Falso | `false` | | Se `true`, verifica os certificados SSL do provedor OpenID. |
| `introspection_endpoint` | string | Falso | | | URL do endpoint de introspecção de token (*token introspection*) do provedor OpenID, usado para introspecção de *access tokens*. Se não estiver definido, o endpoint de introspecção apresentado no documento de descoberta *well-known* será usado como *fallback*. |
| `introspection_endpoint_auth_method` | string | Falso | `client_secret_basic` | | Método de autenticação para o endpoint de introspecção de token. O valor deve ser um dos métodos de autenticação especificados nos metadados `introspection_endpoint_auth_methods_supported` do servidor de autorização, conforme visto no documento de descoberta *well-known*, como `client_secret_basic`, `client_secret_post`, `private_key_jwt` e `client_secret_jwt`. |
| `token_endpoint_auth_method` | string | Falso | `client_secret_basic` | | Método de autenticação para o endpoint de token. O valor deve ser um dos métodos de autenticação especificados nos metadados `token_endpoint_auth_methods_supported` do servidor de autorização, conforme visto no documento de descoberta *well-known*, como `client_secret_basic`, `client_secret_post`, `private_key_jwt` e `client_secret_jwt`. Se o método configurado não for suportado, será feito *fallback* para o primeiro método no array `token_endpoint_auth_methods_supported`. |
| `public_key` | string | Falso | | | Chave pública usada para verificar a assinatura JWT se um algoritmo assimétrico for usado. Fornecer este valor para realizar a verificação do token ignorará a introspecção do token no fluxo de credenciais do cliente. Você pode passar a chave pública no formato `-----BEGIN PUBLIC KEY-----\n……\n-----END PUBLIC KEY-----`. |
| `use_jwks` | boolean | Falso | `false` | | Se `true` e se `public_key` não estiver definida, usa o JWKS para verificar a assinatura JWT e ignora a introspecção do token no fluxo de credenciais do cliente. O endpoint JWKS é analisado a partir do documento de descoberta. |
| `use_pkce` | boolean | Falso | `false` | | Se `true`, usa o Proof Key for Code Exchange (PKCE) para o Fluxo de Código de Autorização, conforme definido no RFC 7636. |
| `token_signing_alg_values_expected` | string | Falso | | | Algoritmo usado para assinar o JWT, como `RS256`. |
| `set_access_token_header` | boolean | Falso | `true` | | Se `true`, define o *access token* em um cabeçalho da solicitação. Por padrão, o cabeçalho `X-Access-Token` é usado. |
| `access_token_in_authorization_header` | boolean | Falso | `false` | | Se `true` e se `set_access_token_header` também for `true`, define o *access token* no cabeçalho `Authorization`. |
| `set_id_token_header` | boolean | Falso | `true` | | Se `true` e se o ID token estiver disponível, define o valor no cabeçalho da solicitação `X-ID-Token`. |
| `set_userinfo_header` | boolean | Falso | `true` | | Se `true` e se os dados de *user info* estiverem disponíveis, define o valor no cabeçalho da solicitação `X-Userinfo`. |
| `set_refresh_token_header` | boolean | Falso | `false` | | Se `true` e se o *refresh token* estiver disponível, define o valor no cabeçalho da solicitação `X-Refresh-Token`. |
| `session` | object | Falso | | | Configuração de sessão usada quando `bearer_only` é `false` e o Plugin usa o fluxo de Código de Autorização. |
| `session.secret` | string | Verdadeiro | | 16 ou mais caracteres | Chave usada para criptografia de sessão e operação HMAC quando `bearer_only` é `false`. |
| `session.cookie` | object | Falso | | | Configurações de *cookie*. |
| `session.cookie.lifetime` | integer | Falso | `3600` | | Tempo de vida do *cookie* em segundos. |
| `session_contents` | object | Falso | | | Configurações do conteúdo da sessão. Se não configurado, todos os dados serão armazenados na sessão. |
| `session_contents.access_token` | boolean | Falso | | | Se `true`, armazena o *access token* na sessão. |
| `session_contents.id_token` | boolean | Falso | | | Se `true`, armazena o ID token na sessão. |
| `session_contents.enc_id_token` | boolean | Falso | | | Se `true`, armazena o ID token criptografado na sessão. |
| `session_contents.user` | boolean | Falso | | | Se `true`, armazena as informações do usuário na sessão. |
| `unauth_action` | string | Falso | `auth` | `["auth","deny","pass"]` | Ação para solicitações não autenticadas. Quando definido como `auth`, redireciona para o endpoint de autenticação do provedor OpenID. Quando definido como `pass`, permite a solicitação sem autenticação. Quando definido como `deny`, retorna respostas `401 unauthenticated` em vez de iniciar o fluxo de concessão de código de autorização. |
| `proxy_opts` | object | Falso | | | Configurações para o servidor proxy atrás do qual o provedor OpenID está. |
| `proxy_opts.http_proxy` | string | Falso | | | Endereço do servidor proxy para solicitações HTTP, por exemplo, `http://<proxy_host>:<proxy_port>`. |
| `proxy_opts.https_proxy` | string | Falso | | | Endereço do servidor proxy para solicitações HTTPS, por exemplo, `http://<proxy_host>:<proxy_port>`. |
| `proxy_opts.http_proxy_authorization` | string | Falso | `Basic [base64 username:password]` | Valor padrão do cabeçalho `Proxy-Authorization` a ser usado com `http_proxy`. Pode ser sobrescrito com um cabeçalho de solicitação `Proxy-Authorization` personalizado. |
| `proxy_opts.https_proxy_authorization` | string | Falso | `Basic [base64 username:password]` | Valor padrão do cabeçalho `Proxy-Authorization` a ser usado com `https_proxy`. Não pode ser sobrescrito com um cabeçalho de solicitação `Proxy-Authorization` personalizado, pois com HTTPS a autorização é concluída ao conectar. |
| `proxy_opts.no_proxy` | string | Falso | | | Lista separada por vírgulas de hosts que não devem ser *proxied*. |
| `authorization_params` | object | Falso | | | Parâmetros adicionais para enviar na solicitação ao endpoint de autorização. |
| `client_rsa_private_key` | string | Falso | | | Chave privada RSA do cliente usada para assinar o JWT para autenticação no OP. Obrigatória quando `token_endpoint_auth_method` é `private_key_jwt`. |
| `client_rsa_private_key_id` | string | Falso | | | ID da chave privada RSA do cliente usada para calcular um JWT assinado. Opcional quando `token_endpoint_auth_method` é `private_key_jwt`. |
| `client_jwt_assertion_expires_in` | integer | Falso | `60` | | Duração de vida do JWT assinado para autenticação no OP, em segundos. Usado quando `token_endpoint_auth_method` é `private_key_jwt` ou `client_secret_jwt`. |
| `renew_access_token_on_expiry` | boolean | Falso | `true` | | Se `true`, tenta renovar silenciosamente o *access token* quando ele expirar ou se um *refresh token* estiver disponível. Se o token falhar ao renovar, redireciona o usuário para reautenticação. |
| `access_token_expires_in` | integer | Falso | | | Tempo de vida do *access token* em segundos se nenhum atributo `expires_in` estiver presente na resposta do endpoint de token. |
| `refresh_session_interval` | integer | Falso | | | Intervalo de tempo para atualizar o ID token do usuário sem exigir reautenticação. Quando não definido, não verificará o tempo de expiração da sessão emitida para o cliente pelo *gateway*. Se definido como `900`, significa atualizar o *id_token* do usuário (ou a sessão no navegador) após 900 segundos sem exigir reautenticação. |
| `iat_slack` | integer | Falso | `120` | | Tolerância do *clock skew* em segundos com a *claim* `iat` em um ID token. |
| `accept_none_alg` | boolean | Falso | `false` | | Defina como `true` se o provedor OpenID não assinar seu ID token, como quando o algoritmo de assinatura está definido como `none`. |
| `accept_unsupported_alg` | boolean | Falso | `true` | | Se `true`, ignora a assinatura do ID token para aceitar algoritmo de assinatura não suportado. |
| `access_token_expires_leeway` | integer | Falso | `0` | | Margem de expiração em segundos para renovação do *access token*. Quando definido como um valor maior que 0, a renovação do token ocorrerá no tempo definido antes da expiração do token. Isso evita erros caso o *access token* expire ao chegar ao servidor de recursos. |
| `force_reauthorize` | boolean | Falso | `false` | | Se `true`, executa o fluxo de autorização mesmo quando um token tiver sido armazenado em *cache*. |
| `use_nonce` | boolean | Falso | `false` | | Se `true`, habilita o parâmetro *nonce* na solicitação de autorização. |
| `revoke_tokens_on_logout` | boolean | Falso | `false` | | Se `true`, notifica ao servidor de autorização que um *refresh* ou *access token* obtido anteriormente não é mais necessário no endpoint de revogação. |
| `jwk_expires_in` | integer | Falso | `86400` | | Tempo de expiração para o *cache* JWK em segundos. |
| `jwt_verification_cache_ignore` | boolean | Falso | `false` | | Se `true`, força a re-verificação para um *bearer token* e ignora quaisquer resultados de verificação em *cache* existentes. |
| `cache_segment` | string | Falso | | | Nome opcional de um segmento de *cache*, usado para separar e diferenciar caches usados por introspecção de token ou verificação JWT. |
| `introspection_interval` | integer | Falso | `0` | | TTL do *access token* introspecionado e em *cache* em segundos. O valor padrão é `0`, o que significa que esta opção não é usada e o Plugin usa por padrão o TTL passado pela *claim* de expiração definida em `introspection_expiry_claim`. Se `introspection_interval` for maior que 0 e menor que o TTL passado pela *claim* de expiração definida em `introspection_expiry_claim`, usa `introspection_interval`. |
| `introspection_expiry_claim` | string | Falso | `exp` | | Nome da *claim* de expiração, que controla o TTL do *access token* introspecionado e em *cache*. |
| `introspection_addon_headers` | array[string] | Falso | | | Usado para anexar valores de cabeçalho adicionais à solicitação HTTP de introspecção. Se o cabeçalho especificado não existir na solicitação original, o valor não será anexado. |
| `claim_validator.issuer.valid_issuers` | string[] | Falso | | | Lista branca (*whitelist*) dos emissores verificados do JWT. Quando não passado pelo usuário, o emissor retornado pelo endpoint de descoberta será usado. Caso ambos estejam ausentes, o emissor não será validado. |
| `claim_schema` | object | Falso | | | Esquema JSON da *claim* de resposta OIDC. Exemplo: `{"type":"object","properties":{"access_token":{"type":"string"}},"required":["access_token"]}` - valida que a resposta contém um campo obrigatório do tipo string `access_token`. |

**NOTA:** `encrypt_fields = {"client_secret"}` também está definido no esquema, o que significa que o campo será armazenado criptografado no etcd. Consulte *encrypted storage fields*.
Além disso, você pode usar Variáveis de Ambiente ou segredos do APISIX para armazenar e referenciar atributos do plugin. O APISIX atualmente suporta o armazenamento de segredos de duas maneiras - Variáveis de Ambiente e HashiCorp Vault.

Por exemplo, use o comando abaixo para definir uma variável de ambiente
```bash
export keycloak_secret=abc
```
e use-a na configuração do plugin conforme abaixo
```json
"client_secret": "$ENV://keycloak_secret"
```

## Exemplos
Os exemplos abaixo demonstram como você pode configurar o Plugin `openid-connect` para diferentes cenários.

> **nota**
> Você pode buscar o `admin_key` do `config.yaml` e salvá-lo em uma variável de ambiente com o seguinte comando:
> ```bash
> admin_key=$(yq '.deployment.admin.admin_key[0].key' conf/config.yaml | sed 's/"//g')
> ```

### Fluxo de Código de Autorização (*Authorization Code Flow*)
O fluxo de código de autorização é definido no RFC 6749, Seção 4.1. Envolve a troca de um código de autorização temporário por um *access token* e é normalmente usado por clientes confidenciais e públicos.

O diagrama a seguir ilustra a interação entre diferentes entidades quando você implementa o fluxo de código de autorização:
[https://static.api7.ai/uploads/2023/11/27/Ga2402sb_oidc-code-auth-flow-revised.png]
(Nota: O conteúdo original da página não incluía o diagrama visual aqui. Ele descreve a sequência.)

Quando uma solicitação de entrada não contém um *access token* em seu cabeçalho nem em um *cookie* de sessão apropriado, o Plugin atua como uma *relying party* e redireciona para o servidor de autorização para continuar o fluxo de código de autorização.

Após a autenticação bem-sucedida, o Plugin mantém o token no *cookie* de sessão, e as solicitações subsequentes usarão o token armazenado no *cookie*.

Consulte *Implement Authorization Code Grant* para um exemplo de como usar o Plugin `openid-connect` para integrar com o Keycloak usando o fluxo de código de autorização.

### Proof Key for Code Exchange (PKCE)
O Proof Key for Code Exchange (PKCE) é definido no RFC 7636. O PKCE melhora o fluxo de código de autorização adicionando um desafio de código (*code challenge*) e um verificador (*verifier*) para prevenir ataques de interceptação de código de autorização.

O diagrama a seguir ilustra a interação entre diferentes entidades quando y (Nota: A captura do conteúdo da página terminou abruptamente aqui, no meio de uma frase sobre PKCE).

See Implement Authorization Code Grant for an example to use the openid-connect Plugin to integrate with Keycloak using the authorization code flow with PKCE.

### Client Credential Flow#
The client credential flow is defined in RFC 6749, Section 4.4. It involves clients requesting an access token with its own credentials to access protected resources, typically used in machine to machine authentication and is not on behalf of a specific user.

The following diagram illustrates the interaction between different entities when you implement the client credential flow:

Client credential flow diagram

See Implement Client Credentials Grant for an example to use the openid-connect Plugin to integrate with Keycloak using the client credentials flow.

### Introspection Flow
The introspection flow is defined in RFC 7662. It involves verifying the validity and details of an access token by querying an authorization server’s introspection endpoint.

In this flow, when a client presents an access token to the resource server, the resource server sends a request to the authorization server’s introspection endpoint, which responds with token details if the token is active, including information like token expiration, associated scopes, and the user or client it belongs to.

The following diagram illustrates the interaction between different entities when you implement the authorization code flow with token introspection:


Client credential with introspection diagram

See Implement Client Credentials Grant for an example to use the openid-connect Plugin to integrate with Keycloak using the client credentials flow with token introspection.

### Password Flow
The password flow is defined in RFC 6749, Section 4.3. It is designed for trusted applications, allowing them to obtain an access token directly using a user’s username and password. In this grant type, the client app sends the user’s credentials along with its own client ID and secret to the authorization server, which then authenticates the user and, if valid, issues an access token.

Though efficient, this flow is intended for highly trusted, first-party applications only, as it requires the app to handle sensitive user credentials directly, posing significant security risks if used in third-party contexts.

The following diagram illustrates the interaction between different entities when you implement the password flow:

Password flow diagram

See Implement Password Grant for an example to use the openid-connect Plugin to integrate with Keycloak using the password flow.

### efresh Token Grant#
The refresh token grant is defined in RFC 6749, Section 6. It enables clients to request a new access token without requiring the user to re-authenticate, using a previously issued refresh token. This flow is typically used when an access token expires, allowing the client to maintain continuous access to resources without user intervention. Refresh tokens are issued along with access tokens in certain OAuth flows and their lifespan and security requirements depend on the authorization server’s configuration.

The following diagram illustrates the interaction between different entities when implementing password flow with refresh token flow:

Password grant with refresh token flow diagram

See Refresh Token for an example to use the openid-connect Plugin to integrate with Keycloak using the password flow with token refreshes.

## Troubleshooting#
This section covers a few commonly seen issues when working with this Plugin to help you troubleshoot.

### APISIX Cannot Connect to OpenID provider#
If APISIX fails to resolve or cannot connect to the OpenID provider, double check the DNS settings in your configuration file config.yaml and modify as needed.

### No Session State Found#
If you encounter a 500 internal server error with the following message in the log when working with authorization code flow, there could be a number of reasons.

the error request to the redirect_uri path, but there's no session state found
1. Misconfigured Redirection URI#
A common misconfiguration is to configure the redirect_uri the same as the URI of the route. When a user initiates a request to visit the protected resource, the request directly hits the redirection URI with no session cookie in the request, which leads to the no session state found error.

To properly configure the redirection URI, make sure that the redirect_uri matches the Route where the Plugin is configured, without being fully identical. For instance, a correct configuration would be to configure uri of the Route to /api/v1/* and the path portion of the redirect_uri to /api/v1/redirect.

You should also ensure that the redirect_uri include the scheme, such as http or https.

2. Cookie Not Sent or Absent#
Check if the SameSite cookie attribute is properly set (i.e. if your application needs to send the cookie cross sites) to see if this could be a factor that prevents the cookie being saved to the browser's cookie jar or being sent from the browser.

3. Upstream Sent Too Big Header#
If you have NGINX sitting in front of APISIX to proxy client traffic, see if you observe the following error in NGINX's error.log:

upstream sent too big header while reading response header from upstream
If so, try adjusting proxy_buffers, proxy_buffer_size, and proxy_busy_buffers_size to larger values.

Another option is to configure the session_content attribute to adjust which data to store in session. For instance, you can set session_content.access_token to true.

4. Invalid Client Secret#
Verify if client_secret is valid and correct. An invalid client_secret would lead to an authentication failure and no token shall be returned and stored in session.

5. PKCE IdP Configuration#
If you are enabling PKCE with the authorization code flow, make sure you have configured the IdP client to use PKCE. For example, in Keycloak, you should configure the PKCE challenge method in the client's advanced settings:

