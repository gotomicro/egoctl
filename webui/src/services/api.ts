import request from "@/utils/request";

export default {
  ProjectList: async (params: any) => {
    return request("/api/projects", {
      method: "GET",
      params,
    });
  },
  ProjectCreate: async (params: any) => {
    return request("/api/projects", {
      method: "POST",
      data: params,
    });
  },
  ProjectUpdate: async (params: any) => {
    return request(`/api/projects/${params.id}`, {
      method: "PUT",
      data: params,
    });
  },
  ProjectDelete: async (params: any) => {
    return request(`/api/projects/${params.id}`, {
      method: "DELETE",
      data: params,
    });
  },
}
