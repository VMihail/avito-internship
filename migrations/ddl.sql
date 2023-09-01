create table if not exists segmentation_service.employee (
    employee_id serial not null primary key
);

create table if not exists segmentation_service.experiment (
    experiment_id serial not null primary key,
    experiment_name varchar not null unique
);

create table if not exists segmentation_service.employee_experiment (
    employee_experiment_id serial not null primary key,
    employee_id serial not null, foreign key (employee_id) references segmentation_service.employee(employee_id),
    experiment_id serial not null, foreign key (experiment_id) references segmentation_service.experiment(experiment_id),
    unique (employee_id, experiment_id)
);

alter table segmentation_service.employee_experiment add if not exists dateAdded date;
alter table segmentation_service.employee_experiment add if not exists dateRemoved date
    check (dateRemoved >= employee_experiment.dateAdded);

alter table segmentation_service.employee_experiment add if not exists date_delete date;

create or replace function delete_employee_experiment_by_timer() returns trigger as $$
begin
    update segmentation_service.employee_experiment
    set dateremoved = date_delete
    WHERE employee_experiment_id = OLD.employee_experiment_id;
    RETURN NULL;
end;
$$ language plpgsql;

create trigger employee_experiment_delete_trigger
    after update on segmentation_service.employee_experiment
    for each row
    when (OLD.date_delete is not null and OLD.date_delete <= NOW())
execute function delete_employee_experiment_by_timer();