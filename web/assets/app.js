const state = {
  scheduleItems: [],
  groups: [],
  subjects: [],
  teachers: [],
  classrooms: []
};

const WEEKDAY_LABELS = {
  1: "Понедельник",
  2: "Вторник",
  3: "Среда",
  4: "Четверг",
  5: "Пятница",
  6: "Суббота",
  7: "Воскресенье"
};

function setStatus(node, message, kind = "muted") {
  node.textContent = message;
  node.classList.remove("ok", "error", "muted");
  node.classList.add(kind);
}

function toInt(value) {
  const parsed = Number.parseInt(value, 10);
  return Number.isFinite(parsed) ? parsed : null;
}

function renderGroupSelects() {
  const selects = [
    document.getElementById("search-group-select"),
    document.getElementById("create-group-select")
  ].filter(Boolean);

  for (const select of selects) {
    select.innerHTML = "";

    const placeholder = document.createElement("option");
    placeholder.value = "";
    placeholder.textContent = state.groups.length ? "Выберите группу" : "Группы не найдены";
    placeholder.disabled = true;
    placeholder.selected = true;
    select.appendChild(placeholder);

    for (const group of state.groups) {
      const option = document.createElement("option");
      option.value = String(group.id);
      option.textContent = group.name;
      select.appendChild(option);
    }
  }
}

function renderDeleteGroupSelect() {
  const select = document.getElementById("delete-group-select");
  if (!select) {
    return;
  }

  select.innerHTML = "";

  const placeholder = document.createElement("option");
  placeholder.value = "";
  placeholder.textContent = state.groups.length ? "Выберите группу" : "Группы не найдены";
  placeholder.disabled = true;
  placeholder.selected = true;
  select.appendChild(placeholder);

  for (const group of state.groups) {
    const option = document.createElement("option");
    option.value = group.name;
    option.textContent = group.name;
    select.appendChild(option);
  }
}

async function loadGroups() {
  const response = await fetch("/groups");
  if (!response.ok) {
    const text = await response.text();
    throw new Error(text || `Ошибка ${response.status}`);
  }

  const payload = await response.json();
  state.groups = Array.isArray(payload) ? payload : [];
  renderGroupSelects();
  renderDeleteGroupSelect();
}

function renderSubjectsSelect() {
  const select = document.getElementById("create-subject-select");
  if (!select) {
    return;
  }

  select.innerHTML = "";

  const placeholder = document.createElement("option");
  placeholder.value = "";
  placeholder.textContent = state.subjects.length ? "Выберите предмет" : "Предметы не найдены";
  placeholder.disabled = true;
  placeholder.selected = true;
  select.appendChild(placeholder);

  for (const subject of state.subjects) {
    const option = document.createElement("option");
    option.value = String(subject.id);
    option.textContent = subject.name;
    select.appendChild(option);
  }
}

function renderTeachersSelect() {
  const select = document.getElementById("create-teacher-select");
  if (!select) {
    return;
  }

  select.innerHTML = "";

  const placeholder = document.createElement("option");
  placeholder.value = "";
  placeholder.textContent = state.teachers.length ? "Выберите преподавателя" : "Преподаватели не найдены";
  placeholder.disabled = true;
  placeholder.selected = true;
  select.appendChild(placeholder);

  for (const teacher of state.teachers) {
    const option = document.createElement("option");
    option.value = String(teacher.id);
    option.textContent = teacher.fullname;
    select.appendChild(option);
  }
}

function renderClassroomsSelect() {
  const select = document.getElementById("create-classroom-select");
  if (!select) {
    return;
  }

  select.innerHTML = "";

  const placeholder = document.createElement("option");
  placeholder.value = "";
  placeholder.textContent = state.classrooms.length ? "Выберите аудиторию" : "Аудитории не найдены";
  placeholder.disabled = true;
  placeholder.selected = true;
  select.appendChild(placeholder);

  for (const classroom of state.classrooms) {
    const option = document.createElement("option");
    option.value = String(classroom.id);
    option.textContent = classroom.number;
    select.appendChild(option);
  }
}

async function loadSubjects() {
  const response = await fetch("/subjects");
  if (!response.ok) {
    const text = await response.text();
    throw new Error(text || `Ошибка ${response.status}`);
  }

  const payload = await response.json();
  state.subjects = Array.isArray(payload) ? payload : [];
  renderSubjectsSelect();
}

async function loadTeachers() {
  const response = await fetch("/teachers");
  if (!response.ok) {
    const text = await response.text();
    throw new Error(text || `Ошибка ${response.status}`);
  }

  const payload = await response.json();
  state.teachers = Array.isArray(payload) ? payload : [];
  renderTeachersSelect();
}

async function loadClassrooms() {
  const response = await fetch("/classrooms");
  if (!response.ok) {
    const text = await response.text();
    throw new Error(text || `Ошибка ${response.status}`);
  }

  const payload = await response.json();
  state.classrooms = Array.isArray(payload) ? payload : [];
  renderClassroomsSelect();
}

async function loadReferenceData() {
  await Promise.all([
    loadGroups(),
    loadSubjects(),
    loadTeachers(),
    loadClassrooms()
  ]);
}

function renderSchedule(items, mode = "pair") {
  const container = document.getElementById("schedule-result");
  const groupNameNode = document.getElementById("result-group-name");
  container.innerHTML = "";
  if (groupNameNode) {
    groupNameNode.textContent = "";
  }

  if (!items.length) {
    container.innerHTML = '<p class="muted">Ничего не найдено по заданным параметрам.</p>';
    return;
  }

  const weekTypeRank = (weekType) => {
    if (weekType === null || weekType === undefined) {
      return 0;
    }
    if (weekType === 1) {
      return 1;
    }
    if (weekType === 2) {
      return 2;
    }
    return 3;
  };

  const sorted = [...items].sort((a, b) => {
    if (a.weekday !== b.weekday) {
      return a.weekday - b.weekday;
    }
    if (a.lesson_number !== b.lesson_number) {
      return a.lesson_number - b.lesson_number;
    }
    if (weekTypeRank(a.week_type) !== weekTypeRank(b.week_type)) {
      return weekTypeRank(a.week_type) - weekTypeRank(b.week_type);
    }
    return String(a.subject_name).localeCompare(String(b.subject_name), "ru");
  });

  const createScheduleItemCard = (item) => {
    const weekTypeLabel = item.week_type === null || item.week_type === undefined
      ? "обе недели"
      : `неделя ${item.week_type}`;

    const block = document.createElement("article");
    block.className = "schedule-item";
    block.innerHTML = `
      <p><strong>${item.subject_name}</strong> (${item.lesson_number} пара, ${weekTypeLabel})</p>
      <p>Преподаватель: ${item.teacher_name}</p>
      <p>Кабинет: ${item.classroom_num}</p>
      <p>Подгруппа: ${item.subgroup ?? "все"}</p>
    `;
    return block;
  };

  if (mode === "week") {
    const groupedByDay = new Map();
    for (const item of sorted) {
      const day = item.weekday;
      if (!groupedByDay.has(day)) {
        groupedByDay.set(day, []);
      }
      groupedByDay.get(day).push(item);
    }

    for (const [day, dayItems] of groupedByDay.entries()) {
      const section = document.createElement("section");
      section.className = "weekday-section";

      const title = document.createElement("h3");
      title.className = "weekday-title";
      title.textContent = WEEKDAY_LABELS[day] || `День ${day}`;
      section.appendChild(title);

      const list = document.createElement("div");
      list.className = "weekday-items";
      for (const item of dayItems) {
        list.appendChild(createScheduleItemCard(item));
      }

      section.appendChild(list);
      container.appendChild(section);
    }
    return;
  }

  for (const item of sorted) {
    container.appendChild(createScheduleItemCard(item));
  }
}

async function fetchWeek(groupID, weekday, weekType, subgroup) {
  const params = new URLSearchParams({
    group_id: String(groupID),
    weekday: String(weekday),
    week_type: String(weekType)
  });

  if (subgroup !== null) {
    params.set("subgroup", String(subgroup));
  }

  const response = await fetch(`/schedule?${params.toString()}`);
  if (!response.ok) {
    const text = await response.text();
    throw new Error(text || `Ошибка ${response.status}`);
  }

  const payload = await response.json();
  if (Array.isArray(payload)) {
    return payload;
  }

  if (Array.isArray(payload.items)) {
    return payload.items;
  }

  return [];
}

async function fetchWeeklySchedule(groupID, weekType, subgroup) {
  const params = new URLSearchParams({
    group_id: String(groupID)
  });

  if (weekType !== null) {
    params.set("week_type", String(weekType));
  }

  if (subgroup !== null) {
    params.set("subgroup", String(subgroup));
  }

  const response = await fetch(`/schedule/week?${params.toString()}`);
  if (!response.ok) {
    const text = await response.text();
    throw new Error(text || `Ошибка ${response.status}`);
  }

  const payload = await response.json();
  if (Array.isArray(payload)) {
    return payload;
  }

  if (Array.isArray(payload.items)) {
    return payload.items;
  }

  return [];
}

function setupSearchForm() {
  const form = document.getElementById("search-form");
  const status = document.getElementById("search-status");
  const weeklyButton = document.getElementById("weekly-schedule-btn");
  const groupNameNode = document.getElementById("result-group-name");

  const renderResultGroupName = (groupID, items) => {
    if (!groupNameNode) {
      return;
    }

    const selectedGroup = state.groups.find((group) => group.id === groupID);
    const groupName = selectedGroup?.name || items[0]?.group_name || "";
    groupNameNode.textContent = groupName ? `Группа: ${groupName}` : "";
  };

  const runSearch = async (mode) => {
    const formData = new FormData(form);
    const groupID = toInt(formData.get("group_id"));
    const weekday = toInt(formData.get("weekday"));
    const lessonNumber = toInt(formData.get("lesson_number"));
    const weekType = toInt(formData.get("week_type"));
    const subgroupRaw = formData.get("subgroup");
    const subgroup = subgroupRaw ? toInt(subgroupRaw) : null;

    if (!groupID) {
      setStatus(status, "Выберите группу.", "error");
      return;
    }

    if (mode === "pair" && (!weekday || !lessonNumber)) {
      setStatus(status, "Проверьте корректность обязательных полей для поиска пары.", "error");
      return;
    }

    setStatus(status, "Запрашиваю данные...", "muted");

    try {
      let resultItems = [];

      if (mode === "week") {
        resultItems = await fetchWeeklySchedule(groupID, weekType, subgroup);
      } else {
        const weeks = weekType ? [weekType] : [1, 2];
        const responses = await Promise.all(weeks.map((w) => fetchWeek(groupID, weekday, w, subgroup)));
        const merged = responses.flat();

        const uniqueById = new Map();
        for (const item of merged) {
          uniqueById.set(item.id, item);
        }

        resultItems = Array.from(uniqueById.values()).filter((item) => item.lesson_number === lessonNumber);
      }

      state.scheduleItems = resultItems;
      renderSchedule(state.scheduleItems, mode);
      renderResultGroupName(groupID, resultItems);

      const modeLabel = mode === "week" ? "на неделю" : "для выбранной пары";
      setStatus(status, `Найдено записей ${modeLabel}: ${resultItems.length}`, "ok");
    } catch (error) {
      setStatus(status, `Ошибка: ${error.message}`, "error");
      renderSchedule([], mode);
      if (groupNameNode) {
        groupNameNode.textContent = "";
      }
    }
  };

  form.addEventListener("submit", async (event) => {
    event.preventDefault();
    await runSearch("pair");
  });

  weeklyButton.addEventListener("click", async () => {
    await runSearch("week");
  });
}

function makePayload(kind, formData) {
  switch (kind) {
    case "teachers":
      return [{ fullname: String(formData.get("fullname")).trim() }];
    case "subjects":
      return [{ Name: String(formData.get("name")).trim() }];
    case "classrooms":
      return [{ Number: String(formData.get("number")).trim() }];
    case "groups":
      return [{ Name: String(formData.get("name")).trim() }];
    default:
      throw new Error("Неизвестный тип формы");
  }
}

function makeDeleteSchedulePayload(formData) {
  const payload = {
    group_name: String(formData.get("group_name")).trim(),
    weekday: toInt(formData.get("weekday")),
    lesson_number: toInt(formData.get("lesson_number"))
  };

  const weektypeRaw = formData.get("weektype");
  if (weektypeRaw !== "") {
    payload.weektype = toInt(weektypeRaw);
  }

  const subgroupRaw = formData.get("subgroup");
  if (subgroupRaw !== "") {
    payload.subgroup = toInt(subgroupRaw);
  }

  return payload;
}

function setupAdminSimpleForms() {
  const forms = document.querySelectorAll(".admin-form");

  forms.forEach((form) => {
    const status = form.querySelector(".status");
    const endpoint = form.dataset.endpoint;
    const kind = form.dataset.kind;

    form.addEventListener("submit", async (event) => {
      event.preventDefault();
      const formData = new FormData(form);

      try {
        const payload = makePayload(kind, formData);
        const response = await fetch(endpoint, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(payload)
        });

        if (!response.ok) {
          const text = await response.text();
          throw new Error(text || `Ошибка ${response.status}`);
        }

        setStatus(status, "Успешно добавлено.", "ok");
        form.reset();

        if (kind === "groups") {
          try {
            await loadGroups();
          } catch (refreshError) {
            setStatus(status, `Добавлено, но не удалось обновить список групп: ${refreshError.message}`, "error");
          }
        }

        if (kind === "subjects") {
          try {
            await loadSubjects();
          } catch (refreshError) {
            setStatus(status, `Добавлено, но не удалось обновить список предметов: ${refreshError.message}`, "error");
          }
        }

        if (kind === "teachers") {
          try {
            await loadTeachers();
          } catch (refreshError) {
            setStatus(status, `Добавлено, но не удалось обновить список преподавателей: ${refreshError.message}`, "error");
          }
        }

        if (kind === "classrooms") {
          try {
            await loadClassrooms();
          } catch (refreshError) {
            setStatus(status, `Добавлено, но не удалось обновить список аудиторий: ${refreshError.message}`, "error");
          }
        }
      } catch (error) {
        setStatus(status, `Ошибка: ${error.message}`, "error");
      }
    });
  });
}

function setupScheduleCreateForm() {
  const form = document.getElementById("schedule-create-form");
  const status = document.getElementById("schedule-create-status");

  form.addEventListener("submit", async (event) => {
    event.preventDefault();
    const formData = new FormData(form);

    const payloadItem = {
      group_id: toInt(formData.get("group_id")),
      subject_id: toInt(formData.get("subject_id")),
      teacher_id: toInt(formData.get("teacher_id")),
      classroom_id: toInt(formData.get("classroom_id")),
      weekday: toInt(formData.get("weekday")),
      lesson_number: toInt(formData.get("lesson_number"))
    };

    const weekType = toInt(formData.get("week_type"));
    if (weekType !== null) {
      payloadItem.week_type = weekType;
    }

    const subgroupRaw = formData.get("subgroup");
    if (subgroupRaw) {
      payloadItem.subgroup = toInt(subgroupRaw);
    }

    const mandatory = [
      payloadItem.group_id,
      payloadItem.subject_id,
      payloadItem.teacher_id,
      payloadItem.classroom_id,
      payloadItem.weekday,
      payloadItem.lesson_number
    ];

    if (mandatory.some((value) => !value)) {
      setStatus(status, "Проверьте, что все обязательные поля заполнены числами.", "error");
      return;
    }

    try {
      const response = await fetch("/schedule", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payloadItem)
      });

      if (!response.ok) {
        const text = await response.text();
        throw new Error(text || `Ошибка ${response.status}`);
      }

      setStatus(status, "Запись в Schedule успешно добавлена.", "ok");
      form.reset();
    } catch (error) {
      setStatus(status, `Ошибка: ${error.message}`, "error");
    }
  });
}

function setupScheduleDeleteForm() {
  const form = document.getElementById("schedule-delete-form");
  const status = document.getElementById("schedule-delete-status");

  form.addEventListener("submit", async (event) => {
    event.preventDefault();
    const formData = new FormData(form);

    const payload = makeDeleteSchedulePayload(formData);

    if (!payload.group_name || !payload.weekday || !payload.lesson_number) {
      setStatus(status, "Проверьте, что группа, день недели и номер пары заполнены корректно.", "error");
      return;
    }

    try {
      const response = await fetch("/schedule", {
        method: "DELETE",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload)
      });

      if (!response.ok) {
        const text = await response.text();
        throw new Error(text || `Ошибка ${response.status}`);
      }

      setStatus(status, "Запись успешно удалена.", "ok");
      form.reset();
    } catch (error) {
      setStatus(status, `Ошибка: ${error.message}`, "error");
    }
  });
}

function setupTabs() {
  const buttons = document.querySelectorAll(".tab-btn");
  const panels = {
    search: document.getElementById("panel-search"),
    admin: document.getElementById("panel-admin")
  };

  buttons.forEach((button) => {
    button.addEventListener("click", () => {
      const tab = button.dataset.tab;

      buttons.forEach((btn) => btn.classList.remove("active"));
      button.classList.add("active");

      Object.entries(panels).forEach(([name, panel]) => {
        panel.classList.toggle("active", name === tab);
      });
    });
  });
}

setupTabs();
setupSearchForm();
setupAdminSimpleForms();
setupScheduleCreateForm();
setupScheduleDeleteForm();
renderSchedule([]);

loadReferenceData().catch((error) => {
  const searchStatus = document.getElementById("search-status");
  setStatus(searchStatus, `Не удалось загрузить справочники: ${error.message}`, "error");
});
