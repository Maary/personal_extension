create table `rms_tasks`
CREATE TABLE IF NOT EXISTS `rms_tasks` (
                                         `id` integer AUTO_INCREMENT NOT NULL PRIMARY KEY,
                                         `page` integer NOT NULL DEFAULT 0 ,
                                         `offset` integer NOT NULL DEFAULT 0 ,
                                         `content` varchar(255) NOT NULL DEFAULT '' ,
  `created` datetime NOT NULL,
  `updated` datetime ,
  `type` varchar(255) NOT NULL DEFAULT '' ,
  `u_u_i_d` varchar(255) NOT NULL DEFAULT ''
  ) ENGINE=InnoDB;