\connect dadosjusbr_test

create table orgaos
(
    id         varchar(10) primary key,
    nome       varchar(100),
    jurisdicao varchar(25),
    entidade   varchar(25),
    uf         varchar(25),
    coletando  json
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


