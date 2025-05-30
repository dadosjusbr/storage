\connect dadosjusbr_test

create table orgaos
(
    id             varchar(10) primary key,
    nome           varchar(100),
    jurisdicao     varchar(25),
    entidade       varchar(25),
    uf             varchar(25),
    coletando      json,
    twitter_handle varchar(25),
    ouvidoria      varchar(100)
);

create table coletas
(
    id                           varchar(25),
    id_orgao                     varchar(10),
    mes                          integer,
    ano                          integer,
    timestamp                    timestamp,
    repositorio_coletor          varchar(150),
    versao_coletor               varchar(150),
    repositorio_parser           varchar(150),
    versao_parser                varchar(150),
    estritamente_tabular         boolean,
    formato_consistente          boolean,
    tem_matricula                boolean,
    tem_lotacao                  boolean,
    tem_cargo                    boolean,
    acesso                       varchar(50),
    extensao                     varchar(25),
    detalhamento_receita_base    varchar(25),
    detalhamento_outras_receitas varchar(25),
    detalhamento_descontos       varchar(25),
    indice_completude            numeric,
    indice_facilidade            numeric,
    indice_transparencia         numeric,
    sumario                      json,
    package                      json,
    procinfo                     json,
    atual                        boolean,
    backups                      json,
    formato_aberto               boolean,
    duracao_segundos             double precision,
    manual                       boolean,

    constraint coleta_pk primary key (id,timestamp),
    constraint coleta_orgao_fk foreign key (id_orgao) references orgaos(id) on delete cascade
);

create table remuneracoes_zips
(
    id_orgao         varchar(10),
    mes              integer,
    ano              integer,
    linhas_descontos integer,
    linhas_base      integer,
    linhas_outras    integer,
    zip_url          text,

    constraint remuneracoes_pk primary key (id_orgao, mes, ano )
);

create table contracheques
(
    id integer,
    orgao varchar(10),
    mes integer,
    ano integer,
    chave_coleta varchar(20),
    nome varchar(100),
    matricula varchar(20),
    funcao varchar(100),
    local_trabalho varchar(100),
    salario numeric,
    beneficios numeric,
    descontos numeric,
    remuneracao numeric,
    situacao varchar(5),
    nome_sanitizado varchar(150),

    constraint contracheques_pk primary key (id, orgao, mes, ano)
);

create table remuneracoes
(
    id integer,
    id_contracheque integer,
    orgao varchar(10),
    mes integer,
    ano integer,
    categoria varchar(100),
    item varchar(100),
    valor numeric,
    inconsistente boolean,
    tipo varchar(5),
    item_sanitizado varchar(100),

    constraint pk_remuneracoes primary key (id, id_contracheque, orgao, mes, ano),
    constraint fk_remuneracoes foreign key (id_contracheque, orgao, mes, ano) references contracheques(id, orgao, mes, ano) on delete cascade
);

CREATE MATERIALIZED VIEW public.media_por_membro
TABLESPACE pg_default
AS SELECT media_por_membro.orgao,
    media_por_membro.ano,
    avg(media_por_membro.salario) AS salario,
    avg(media_por_membro.beneficios) AS beneficios,
    avg(media_por_membro.descontos) AS descontos,
    avg(media_por_membro.remuneracao) AS remuneracao
   FROM ( SELECT c.orgao,
            c.ano,
            c.nome_sanitizado,
            count(*) AS num_meses,
            avg(c.salario) AS salario,
            avg(c.beneficios) AS beneficios,
            avg(c.descontos) AS descontos,
            avg(c.remuneracao) AS remuneracao
           FROM contracheques c
          GROUP BY c.orgao, c.ano, c.nome_sanitizado) media_por_membro
  WHERE media_por_membro.num_meses > 1
  GROUP BY media_por_membro.orgao, media_por_membro.ano
WITH DATA;

CREATE MATERIALIZED VIEW orgao_mes_ano_inconsistentes AS 
WITH remuneracoes_inconsistentes AS (
			SELECT DISTINCT orgao, ano, mes
			FROM remuneracoes
			WHERE inconsistente = TRUE
			)
			SELECT c.id_orgao, c.ano, c.mes, 
				CASE 
					WHEN ri.orgao IS NOT NULL THEN TRUE 
					ELSE FALSE 
				END AS inconsistente
			FROM coletas c
			LEFT JOIN remuneracoes_inconsistentes ri
				ON ri.orgao = c.id_orgao
				AND ri.ano = c.ano
				AND ri.mes = c.mes
			WHERE c.atual = TRUE
			AND (c.procinfo::text = 'null' OR c.procinfo IS NULL);

CREATE MATERIALIZED VIEW orgao_ano_inconsistentes AS 
WITH remuneracoes_inconsistentes AS (
			SELECT DISTINCT orgao, ano
			FROM remuneracoes
			WHERE inconsistente = TRUE
			)
			SELECT distinct c.id_orgao, c.ano, 
				CASE 
					WHEN ri.orgao IS NOT NULL THEN TRUE 
					ELSE FALSE 
				END AS inconsistente
			FROM coletas c
			LEFT JOIN remuneracoes_inconsistentes ri
				ON ri.orgao = c.id_orgao
				AND ri.ano = c.ano
			WHERE c.atual = true
			AND (c.procinfo::text = 'null' OR c.procinfo IS NULL);