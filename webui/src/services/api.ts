import request from "@/utils/request";

export default {
  ProjectList: async (params: any) => {
    return request("/api/projects", {
      method: "GET",
      params,
    });
  },
  ProjectGen: async (params: any) => {
    return request(`/api/projects/gen`, {
      method: "GET",
      params: {
        path: params.path,
      },
    });
  },
  ProjectCreate: async (params: any) => {
    return request("/api/projects", {
      method: "POST",
      data: params,
    });
  },
  ProjectUpdate: async (params: any) => {
    return request(`/api/projects`, {
      method: "PUT",
      data: params,
    });
  },
  ProjectUpdateDSL: async (params: any) => {
    return request(`/api/projects/dsl`, {
      method: "PUT",
      data: params,
    });
  },
  ProjectDelete: async (params: any) => {
    return request(`/api/projects`, {
      method: "DELETE",
      data: params,
    });
  },
  TemplateList: async (params: any) => {
    return request("/api/templates", {
      method: "GET",
      params,
    });
  },
  TemplateSelect: async (params: any) => {
    return request("/api/templates/select", {
      method: "GET",
      params,
    });
  },
  TemplateCreate: async (params: any) => {
    return request("/api/templates", {
      method: "POST",
      data: params,
    });
  },
  TemplateUpdate: async (params: any) => {
    return request(`/api/templates`, {
      method: "PUT",
      data: params,
    });
  },
  TemplateSync: async (params: any) => {
    return request(`/api/templates/sync`, {
      method: "PUT",
      data: params,
    });
  },
  TemplateDelete: async (params: any) => {
    return request(`/api/templates`, {
      method: "DELETE",
      data: params,
    });
  },
}
